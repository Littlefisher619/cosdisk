package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	//"time"

	"github.com/sirupsen/logrus"
	cosapi "github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

type Cos struct {
	c         *cosapi.Client
	l         *logrus.Entry
	secretID  string
	secretKey string
}

func New(cosurl string, secretID string, secretKey string, logger *logrus.Entry) *Cos {
	logger.Info("New Cos: ", cosurl)
	u, _ := url.Parse(cosurl)
	b := &cosapi.BaseURL{BucketURL: u}
	return &Cos{
		cosapi.NewClient(b, &http.Client{
			Transport: &cosapi.AuthorizationTransport{
				SecretID:  secretID,
				SecretKey: secretKey,
				Transport: &debug.DebugRequestTransport{
					RequestHeader:  true,
					RequestBody:    false,
					ResponseHeader: true,
					ResponseBody:   false,
				},
			},
		}),
		logger,
		secretID,
		secretKey,
	}
}

func (s *Cos) isFileExist(ctx context.Context, filename string) (bool, error) {
	isexist, err := s.c.Object.IsExist(context.Background(), filename)
	if err != nil {
		s.l.Errorf("Fail to check exist %s: %s", filename, err)
		return false, err
	}
	return isexist, nil
}

func (s *Cos) UploadByPath(ctx context.Context, filename string, localpath string) error {
	exist, err := s.isFileExist(ctx, filename)
	if exist {
		return nil
	}

	if err != nil {
		return fmt.Errorf("check exist: %s", filename)
	}

	_, err = s.c.Object.PutFromFile(context.Background(), filename, localpath, nil)
	if err != nil {
		s.l.Errorf("Fail to PutFromFile %s: %s", filename, err)
		return fmt.Errorf("upload file failed: %s", filename)
	}

	return nil
}

func (s *Cos) UploadByReader(ctx context.Context, filename string, io io.Reader) error {
	exist, err := s.isFileExist(ctx, filename)
	if exist {
		return nil
	}

	if err != nil {
		return fmt.Errorf("existence check failed: %s", filename)
	}

	_, err = s.c.Object.Put(ctx, filename, io, nil)
	if err != nil {
		s.l.Errorf("Fail to Upload file %s: %s", filename, err)
		return fmt.Errorf("upload file failed: %s", filename)
	}
	return nil
}

func (s *Cos) Download(ctx context.Context, filename string, localpath string) error {
	_, err := s.c.Object.GetToFile(ctx, filename, localpath, nil)
	if err != nil {
		s.l.Errorf("Fail to GetToFile %s: %s", filename, err)
		return fmt.Errorf("download file failed: %s", filename)
	}
	return nil
}

func (s *Cos) GetReader(ctx context.Context, filename string) (io.Reader, error) {
	resp, err := s.c.Object.Get(ctx, filename, nil)
	if err != nil {
		s.l.Errorf("Fail to Get %s: %s", filename, err)
		return nil, fmt.Errorf("download file failed: %s", filename)
	}
	return resp.Body, nil
}

func (s *Cos) GetDownloadURL(ctx context.Context, filename string) (string, error) {
	url, err := s.c.Object.GetPresignedURL(ctx, http.MethodGet, filename, s.secretID, s.secretKey, time.Minute*1, nil)
	if err != nil {
		s.l.Errorf("Fail to GetPresignedURL for %s: %s", filename, err)
		return "", fmt.Errorf("get url failed: %s", filename)
	}

	return url.String(), nil
}

func (s *Cos) Delete(ctx context.Context, filename string) error {
	_, err := s.c.Object.Delete(ctx, filename, nil)
	if err != nil {
		s.l.Errorf("Fail to Delete %s: %s", filename, err)
		return err
	}
	return nil
}
