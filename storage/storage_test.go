package storage

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

var (
	cosURL = "https://aaa-1251109145.cos.ap-nanjing.myqcloud.com"
	logger *logrus.Entry
)

func init() {
	logger = logrus.New().WithField("subject", "logger on test")
}

func TestByPathandReader(t *testing.T) {
	ass := require.New(t)
	s := New(cosURL, "", "", logger)
	testFile := "storage.go"
	fp, err := os.Open(testFile)
	ass.NoError(err)
	defer fp.Close()
	localcontent, err := ioutil.ReadAll(fp)
	ass.NoError(err)

	err = s.UploadByPath("test.go", testFile)
	ass.NoError(err, "Upload fail")
	// update same file again
	err = s.UploadByPath("test.go", testFile)
	ass.NoError(err, "Upload fail")

	// update by reader
	fp1, err := os.Open(testFile)
	ass.NoError(err)
	defer fp1.Close()
	err = s.UploadByReader(context.Background(), "test1.go", fp1)
	ass.NoError(err, "Upload fail")
	err = s.UploadByReader(context.Background(), "test1.go", fp1)
	ass.NoError(err, "Upload fail")

	err = s.DownloadByPath("test.go", "test.go")
	ass.NoError(err)
	err = s.DownloadByPath("test2.go", "test.go")
	ass.Error(err)

	reader, err := s.DownloadByReader(context.Background(), "test1.go")
	ass.NoError(err)
	_, err = s.DownloadByReader(context.Background(), "test3.go")
	ass.Error(err)

	err = s.Delete("test.go")
	ass.NoError(err)
	err = s.Delete("test1.go")
	ass.NoError(err)

	rc, err := ioutil.ReadAll(reader)
	ass.NoError(err)
	ass.Equal(rc, localcontent)

	testfp, err := os.Open("test.go")
	ass.NoError(err)
	content, err := ioutil.ReadAll(testfp)
	ass.NoError(err)
	ass.Equal(content, localcontent)
	os.Remove("test.go")
}

func TestFakeurl(t *testing.T) {
	ass := require.New(t)

	s := New("fakeurl", "", "", logger)
	err := s.UploadByPath("test.go", "cosc_test.go")
	ass.Error(err)

	err = s.DownloadByPath("test.go", "cosc_test.go")
	ass.Error(err)
}

/*
func TestDownloadByURL(t *testing.T) {
	ass := require.New(t)
	s := New(cosURL, "", "", logger)

	err := s.UploadByPath("test.go", "client.go")
	ass.NoError(err, "Upload fail")

	url, err := s.DownloadByUrl(context.Background(), "test.go")
	ass.NoError(err)
	resp, err := http.Get(url)
	ass.NoError(err)
	defer resp.Body.Close()
	logger.Info(url)
	logger.Info(resp)

	testfp, err := os.Open("client.go")
	ass.NoError(err)
	localcontent, err := ioutil.ReadAll(testfp)
	ass.NoError(err)
	content, err := ioutil.ReadAll(resp.Body)
	ass.NoError(err)
	ass.Equal(content, localcontent)

	err = s.Delete("test.go")
	ass.NoError(err)
}
*/
