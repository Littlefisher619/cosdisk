package userfile

import (
	"fmt"
	"os"
	ospath "path"
	"strings"

	"github.com/Littlefisher619/cosdisk/model"
)

type UserfileStorageKV struct {
	userfileKV model.KeyValueStorage
}

type UserfileTxnKV struct {
	txn model.KeyValueTXN
}

func NewUserfileKV(kv model.KeyValueStorage) *UserfileStorageKV {
	return &UserfileStorageKV{
		userfileKV: kv,
	}
}

func (us *UserfileStorageKV) StartTranscation() (txn model.UserfileTXN, err error) {
	tn, err := us.userfileKV.StartTranscation()
	if err != nil {
		return nil, err
	}
	return &UserfileTxnKV{tn}, nil
}

func (us *UserfileStorageKV) RunInTranscation(txnfunc func(txn model.UserfileTXN) error) error {
	txn, err := us.userfileKV.StartTranscation()
	if err != nil {
		return err
	}

	err = txnfunc(&UserfileTxnKV{txn})
	if err != nil {
		if e := txn.RollingBackTranscation(); e != nil {
			return fmt.Errorf("Rollback failed %s: %w", err, e)
		}
		return err
	}

	err = txn.CommitTranscation()
	if err != nil {
		if e := txn.RollingBackTranscation(); e != nil {
			return e
		}
		return err
	}
	return nil
}

func (us *UserfileTxnKV) CommitTranscation() (err error) {
	return us.txn.CommitTranscation()
}

func (us *UserfileTxnKV) RollingBackTranscation() (err error) {
	return us.txn.RollingBackTranscation()
}

func (us *UserfileTxnKV) ListFiles(userid string, dirpath string) (files []os.FileInfo, err error) {
	if dirpath[len(dirpath)-1] != '/' {
		dirpath = dirpath + "/"
	}
	key := GenerateUserFileKey(userid, dirpath)
	value, err := us.txn.Get(key)
	if err != nil {
		return nil, model.ErrNotFound
	}
	kind, data, err := ParseUserFileValue(value)
	if err != nil {
		return nil, err
	}
	if kind != DirType {
		return nil, model.ErrIsFile
	}
	return ParseDirContent(data)
}

func (us *UserfileTxnKV) GetFileInfo(username string, path string) (info os.FileInfo, err error) {
	key := GenerateUserFileKey(username, path)
	value, err := us.txn.Get(key)
	if err != nil {
		if path[len(path)-1] != '/' {
			path = path + "/"
		} else {
			path = path[:len(path)-1]
		}
		key = GenerateUserFileKey(username, path)
		value, err = us.txn.Get(key)
		if err != nil {
			return nil, model.ErrNotFound
		}
	}
	kind, _, err := ParseUserFileValue(value)
	if err != nil {
		return nil, err
	}
	if kind == FileType {
		return &FileData{
			FName:  ospath.Base(path),
			FIsDir: false,
		}, nil
	}
	return &FileData{
		FName:  ospath.Base(path),
		FIsDir: true,
		FMode:  os.ModeDir,
	}, nil
}

func (us *UserfileTxnKV) GetFileID(userid string, path string) (id string, err error) {
	key := GenerateUserFileKey(userid, path)
	value, err := us.txn.Get(key)
	if err != nil {
		return "", model.ErrNotFound
	}
	kind, data, err := ParseUserFileValue(value)
	if err != nil {
		return "", err
	}
	if kind != FileType {
		return "", model.ErrUnsupportDirTarget
	}
	return data[0], nil
}

func (us *UserfileTxnKV) DeleteFile(userid string, filepath string) error {
	key := GenerateUserFileKey(userid, filepath)
	value, err := us.txn.Get(key)
	if err != nil {
		return model.ErrNotFound
	}
	if strings.HasPrefix(value, "dir") {
		return model.ErrUnsupportDirTarget
	}
	err = us.txn.Delete(key)
	if err != nil {
		return err
	}
	parentpath, filename := Makefilepath(filepath)
	return us.removeFilefromDir(userid, parentpath, filename)
}

func (us *UserfileTxnKV) DeleteDir(userid string, filepath string) error {
	if filepath[len(filepath)-1] != '/' {
		filepath = filepath + "/"
	}
	key := GenerateUserFileKey(userid, filepath)
	value, err := us.txn.Get(key)
	if err != nil {
		return model.ErrNotFound
	}
	if strings.HasPrefix(value, "file") {
		return model.ErrIsFile
	}
	err = us.txn.Delete(key)
	if err != nil {
		return err
	}
	if filepath == "/" {
		return nil
	}
	parentpath, filename := Makedirpath(filepath)
	return us.removeFilefromDir(userid, parentpath, filename)
}

func (us *UserfileTxnKV) AddFile(userid string, filepath string, id string) error {
	if info, err := us.GetFileInfo(userid, filepath); err == nil {
		if info.IsDir() {
			return model.ErrUnsupportDirTarget
		}
		if filepath[len(filepath)-1] != '/' {
			if info, err := us.GetFileInfo(userid, filepath+"/"); err == nil && info.IsDir() {
				return model.ErrUnsupportDirTarget
			}
		}
		anotherid, err := us.GetFileID(userid, filepath)
		if err != nil {
			return err
		}
		if id == anotherid {
			// don't need to do anything
			return nil
		}
		// update file
		err = us.DeleteFile(userid, filepath)
		if err != nil {
			return err
		}
		return us.AddFile(userid, filepath, id)
	}
	parentpath, filename := Makefilepath(filepath)
	err := us.addFileToDir(userid, parentpath, filename)
	if err != nil {
		return err
	}
	err = us.txn.Set(GenerateUserFileKey(userid, filepath), GenerateUserFileValue(id))
	return err
}

func (us *UserfileTxnKV) AddDir(userid string, dirpath string) error {
	if dirpath != "/" {
		// is not root
		info, err := us.GetFileInfo(userid, dirpath)
		if err == nil {
			if !info.IsDir() {
				return model.ErrIsFile
			}
			return model.ErrFileExists
		}

		if err != model.ErrNotFound {
			return fmt.Errorf("fail to create dir: %s", err)
		}

		if dirpath[len(dirpath)-1] != '/' {
			dirpath = dirpath + "/"
			if _, err := us.GetFileInfo(userid, dirpath); err == nil {
				return model.ErrFileExists
			}

			if err != model.ErrNotFound {
				return fmt.Errorf("fail to create dir: %s", err)
			}

		}
		parentpath, dirname := Makedirpath(dirpath)
		err := us.addFileToDir(userid, parentpath, dirname)
		if err != nil {
			return err
		}
	} else if _, err := us.GetFileInfo(userid, dirpath); err == nil {
		return model.ErrFileExists
	}
	err := us.txn.Set(GenerateUserFileKey(userid, dirpath), GenerateUserDirValue(nil))
	if err != nil {
		return err
	}
	return nil
}

func (us *UserfileTxnKV) addFileToDir(userid string, dirpath string, filename string) error {
	key := GenerateUserFileKey(userid, dirpath)
	dirContent, err := us.txn.Get(key)
	if err != nil {
		return model.ErrNotFound
	}
	kind, data, err := ParseUserFileValue(dirContent)
	if err != nil {
		return err
	}
	if kind != DirType {
		return model.ErrIsFile
	}
	data = append(data, filename)
	err = us.txn.Set(key, GenerateUserDirValue(data))
	if err != nil {
		return err
	}
	return nil
}

func (us *UserfileTxnKV) removeFilefromDir(userid string, dirpath string, filename string) error {
	key := GenerateUserFileKey(userid, dirpath)
	dirContent, err := us.txn.Get(key)
	if err != nil {
		return model.ErrNotFound
	}
	kind, data, err := ParseUserFileValue(dirContent)
	if err != nil {
		return err
	}
	if kind != DirType {
		return model.ErrIsFile
	}
	if len(data) == 0 {
		return nil
	}
	var index = -1
	for i, v := range data {
		if v == filename {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	data = append(data[:index], data[index+1:]...)
	err = us.txn.Set(key, GenerateUserDirValue(data))
	if err != nil {
		return err
	}
	return nil
}
