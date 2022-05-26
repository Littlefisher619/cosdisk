package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/Littlefisher619/cosdisk/config"
	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/storage"

	"github.com/sirupsen/logrus"
)

const (
	cosURL = "https://cosdisk-1251109145.cos.ap-shanghai.myqcloud.com"
)

type CosDisk struct {
	Logger              *logrus.Entry
	storage             *storage.Cos
	accountRepository   model.AccountRepository
	userfileRepository  model.UserfileRepository
	shareFileRepository model.ShareFileRepository
}

func New(
	config *config.Config,
) *CosDisk {
	logger := logrus.New().WithField("serviceName", "cosdisk")
	return &CosDisk{
		accountRepository:   config.AccountRepository,
		userfileRepository:  config.UserfileRepository,
		shareFileRepository: config.ShareFileRepository,
		storage:             storage.New(config.CosURL, config.SecretID, config.SecretKey, logger.WithField("layer", "cos")),
		Logger:              logger,
	}
}

func (c *CosDisk) UploadUserFileByReader(ctx context.Context, userId string, filePath string, reader io.Reader) error {
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	logger := c.Logger.WithField("User", userId).WithField("Op", ("UploadUserFileByReader "))
	//logger.Info(filePath + " " + string(buffer))

	reader = bytes.NewReader(buffer)
	h := sha256.New()
	_, err = h.Write(buffer)
	if err != nil {
		return err
	}
	fileId := hex.EncodeToString(h.Sum(nil))
	c.Logger.Info("field id: " + fileId)

	err = c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		err = txn.AddFile(userId, filePath, fileId)
		if err != nil {
			logger.Error(("Add file failed "), err)
			return err
		}

		err = c.storage.UploadByReader(ctx, fileId, reader)
		if err != nil {
			logger.Error(("Upload file failed "), err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *CosDisk) CreateDir(userId string, Path string) error {
	logger := c.Logger.WithField("User", userId).WithField("Op", ("CreateDir "))
	logger.Info(Path)

	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		return txn.AddDir(userId, Path)
	})
	if err != nil {
		logger.Error(("AddDir failed "), err)
		return err
	}

	return nil
}

func (c *CosDisk) DeleteDir(userId string, Path string) error {
	logger := c.Logger.WithField("User", userId).WithField("Op", ("DeleteDir "))
	logger.Info(Path)

	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		return txn.DeleteDir(userId, Path)
	})
	if err != nil {
		logger.Error("DeleteDir failed")
		return err
	}
	return nil
}

func (c *CosDisk) DeleteFIle(userId string, Path string) error {
	logger := c.Logger.WithField("User", userId).WithField("Op", ("DeleteFIle "))
	logger.Info(Path)

	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		return txn.DeleteFile(userId, Path)
	})
	if err != nil {
		logger.Error(("DeleteFile failed "), err)
		return err
	}
	return nil
}

func (c *CosDisk) GetFileInfo(userId string, Path string) (os.FileInfo, error) {
	logger := c.Logger.WithField("User", userId).WithField("Op", "GetFileInfo ")
	logger.Info(Path)
	var info os.FileInfo
	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		var err error
		info, err = txn.GetFileInfo(userId, Path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(("GetFileInfo failed "), err)
		return nil, err
	}
	return info, nil
}

func (c *CosDisk) ListFiles(userId string, Path string) ([]os.FileInfo, error) {
	logger := c.Logger.WithField("User", userId).WithField("Op", "ListFiles ")
	logger.Info(Path)
	var info []os.FileInfo
	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		var err error
		info, err = txn.ListFiles(userId, Path)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error("ListFiles failed ", err)
		return nil, err
	}
	return info, nil
}

func (c *CosDisk) UploadUserFileByPath(userId string, fileName string, filePath string) error {
	return model.ErrUnimplemented
}

func (c *CosDisk) DownloadUserFileByUrl(ctx context.Context, userId string, path string) (string, error) {
	logger := c.Logger.WithField("User", userId).WithField("Op", "DownloadUserFileByUrl ")

	var url string
	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		var err error
		id, err := txn.GetFileID(userId, path)
		if err != nil {
			return err
		}
		url, err = c.storage.GetDownloadURL(ctx, id)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(("DownloadUserFileByUrl failed "), err)
		return "", err
	}
	return url, nil
}

func (c *CosDisk) MoveFile(ctx context.Context, userId string, path string, newpath string) (string, error) {
	c.Logger.Info("MoveFile " + path + " " + newpath)
	logger := c.Logger.WithField("User", userId).WithField("Op", "MoveFile ")
	logger.Info(path + " " + newpath)

	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		var err error
		id, err := txn.GetFileID(userId, path)
		if err != nil {
			return err
		}
		err = txn.AddFile(userId, newpath, id)
		if err != nil {
			logger.Error("Add file failed ", err)
			return err
		}
		err = txn.DeleteFile(userId, path)
		if err != nil {
			logger.Error("Delete file failed ", err)
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(("MoveFile failed "), err)
		return "", err
	}
	return "success", nil
}

func (c *CosDisk) CreateShareFile(ctx context.Context, userId string, path string, expireDays int) (string, error) {
	logger := c.Logger.WithField("User", userId).WithField("Op", "CreateShareFileUrl ")
	logger.Info(path)
	var shareID string
	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		var err error
		_, err = txn.GetFileID(userId, path)
		if err != nil {
			return err
		}
		shareID, err = c.shareFileRepository.CreateShareFile(ctx, userId, path, time.Duration(expireDays)*24*time.Hour)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(("CreateShareFileUrl failed "), err)
		return "", err
	}
	return shareID, err
}

func (c *CosDisk) ShareFileToUser(ctx context.Context, userId string, shareId string) error {
	logger := c.Logger.WithField("User", userId).WithField("Op", "ShareFileToUser ")
	logger.Info(shareId)

	share, err := c.shareFileRepository.GetVaildShareFile(ctx, shareId)
	if err != nil {
		return err
	}
	err = c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		fileID, err := txn.GetFileID(share.UserId, share.Path)
		if err != nil {
			return err
		}
		fileInfo, err := txn.GetFileInfo(share.UserId, share.Path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			logger.Error(("is dir share failed "), err)
			return model.ErrUnimplemented
		}
		// fix me: add file to another dir
		err = txn.AddFile(userId, "/"+fileInfo.Name(), fileID)
		if err != nil {
			logger.Error(("AddFile in share failed "), err)
			return err
		}
		return nil
	})
	if err != nil {
		logger.Error(("ShareFileToUser failed "), err)
		return err
	}
	return nil
}

func (c *CosDisk) DownloadUserFileByReader(ctx context.Context, userId string, path string) (io.Reader, error) {
	logger := c.Logger.WithField("User", userId).WithField("Op", "DownloadUserFileByReader ")
	logger.Info(path)

	var reader io.Reader
	err := c.userfileRepository.RunInTranscation(func(txn model.UserfileTXN) error {
		var err error
		id, err := txn.GetFileID(userId, path)
		if err != nil {
			return err
		}
		reader, err = c.storage.GetReader(ctx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		logger.Error("DownloadUserFileByReader failed ", err)
		return nil, err
	}
	return reader, nil
}
