package storage

import (
	"context"
	"io"
	"os"
	"regexp"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
)

type CosAws struct {
	c         *s3.S3
	l         *logrus.Entry
	secretID  string
	secretKey string
	bucket    string
}

func newSession(cosurl string, secretID string, secretKey string) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(secretID, secretKey, "")
	endpoint := cosurl
	r := regexp.MustCompile(`^https://aaa-1251109145.cos.ap-(.*).myqcloud.com$`)
	region := r.FindStringSubmatch(cosurl)

	config := &aws.Config{
		Region:           aws.String(region[0]),
		Endpoint:         &endpoint,
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
		//DisableSSL:       &disableSSL,
	}
	return session.NewSession(config)
}

func NewAws(cosurl string, secretID string, secretKey string, logger *logrus.Entry) *CosAws {
	logger.Info("New Cos: ", cosurl)
	s, err := newSession(cosurl, secretID, secretKey)
	if err != nil {
		logger.Errorf("New Cos: %s failed %s", cosurl, err)
		return nil
	}
	r := regexp.MustCompile(`^https://(.*).cos.ap-(.*).myqcloud.com$`)
	region := r.FindStringSubmatch(cosurl)
	return &CosAws{
		s3.New(s),
		logger,
		secretID,
		secretKey,
		region[0],
	}
}

func (s *CosAws) isFileExist(filename string) bool {
	/*
		isexist, err := s.c.WaitUntilBucketExists(context.Background(), filename)
		if err != nil {
			s.l.Error("Check ", filename, " is exist failed", " ", err)
			return false
		}
	*/
	// TODO
	return false
}

func (s *CosAws) UploadByReader(ctx context.Context, filename string, io io.Reader) error {
	if s.isFileExist(filename) {
		return nil
	}
	_, err := s.c.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   os.Stdin,
	})
	if err != nil {
		s.l.Error("Upload ", filename, " is exist failed", " ", err)
	}
	return err
}

func (s *CosAws) DownloadByReader(ctx context.Context, filename string) (io.Reader, error) {
	// resp, err := s.c.Object.Get(context.Background(), filename, nil)
	// if err != nil {
	// 	s.l.Error("Download ", filename, " is exist failed", " ", err)
	// }
	// TODO
	return nil, model.ErrUnimplemented
}

func (s *CosAws) DownloadByUrl(ctx context.Context, filename string) (string, error) {
	//url, err := s.c.Object.GetPresignedURL(ctx, http.MethodGet, filename, s.secretID, s.secretKey, time.Hour, nil)
	//url := s.c.(filename)
	return "", model.ErrUnimplemented
}

func (s *CosAws) Delete(filename string) error {
	_, err := s.c.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(filename)})
	if err != nil {
		s.l.Errorf("Unable to delete object %q from bucket %q, %v", filename, s.bucket, err)
	}
	err = s.c.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		s.l.Error("Delete ", filename, " is exist failed", " ", err)
	}
	return err
}
