package storage

import (
	"context"
	"io"
	"net/http"
	"net/url"

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

func (s *Cos) isFileExist(filename string) bool {
	isexist, err := s.c.Object.IsExist(context.Background(), filename)
	if err != nil {
		s.l.Error("Check ", filename, " is exist failed", " ", err)
		return false
	}
	return isexist
}

func (s *Cos) UploadByPath(filename string, localpath string) error {
	if s.isFileExist(filename) {
		return nil
	}
	_, err := s.c.Object.PutFromFile(context.Background(), filename, localpath, nil)
	if err != nil {
		s.l.Error("Upload ", filename, " is exist failed", " ", err)
	}
	return err
}

func (s *Cos) UploadByReader(ctx context.Context, filename string, io io.Reader) error {
	if s.isFileExist(filename) {
		return nil
	}
	_, err := s.c.Object.Put(ctx, filename, io, nil)
	if err != nil {
		s.l.Error("Upload ", filename, " is exist failed", " ", err)
	}
	return err
}

func (s *Cos) DownloadByPath(filename string, localpath string) error {
	_, err := s.c.Object.GetToFile(context.Background(), filename, localpath, nil)
	if err != nil {
		s.l.Error("Download ", filename, " is exist failed", " ", err)
	}
	return err
}

func (s *Cos) DownloadByReader(ctx context.Context, filename string) (io.Reader, error) {
	resp, err := s.c.Object.Get(context.Background(), filename, nil)
	if err != nil {
		s.l.Error("Download ", filename, " is exist failed", " ", err)
	}
	return resp.Body, err
}

func (s *Cos) DownloadByUrl(ctx context.Context, filename string) (string, error) {
	//url, err := s.c.Object.GetPresignedURL(ctx, http.MethodGet, filename, s.secretID, s.secretKey, time.Hour, nil)
	url := s.c.Object.GetObjectURL(filename)
	/*
		if err != nil {
			s.l.Error("create Download url for %s failed", filename)
			return "", err
		}
	*/
	return url.String(), nil
}

func (s *Cos) Delete(filename string) error {
	_, err := s.c.Object.Delete(context.Background(), filename, nil)
	if err != nil {
		s.l.Error("Delete ", filename, " is exist failed", " ", err)
	}
	return err
}
