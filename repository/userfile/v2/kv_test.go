package userfile

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	model "github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/pkg/dbdriver"
)

var client model.UserfileRepository

func init() {
	rc := dbdriver.InitMapKV()
	if rc == nil {
		return
	}
	client = NewClient(rc)
}

func TestAddFile(t *testing.T) {
	assert := require.New(t)
	defer client.Close()
	txn, err := client.Begin()

	assert.NoError(err)
	defer txn.Rollback()
	ctx := context.Background()
	userId := "test"

	err = txn.CreateRepository(ctx, userId)
	assert.NoError(err)

	var rootInfo *model.FileInfoWrapper
	rootInfo, err = txn.GetFileInfo(ctx, userId, "/")
	assert.NoError(err)
	assert.Nil(rootInfo.DirContent)
	assert.Nil(rootInfo.FileObject)
	assert.Equal(true, rootInfo.FIsDir)
	assert.Equal(os.ModeDir, rootInfo.FMode)
	assert.Equal("/", rootInfo.FName)

	testFile1 := &model.FileDetail{
		FileInfo: model.FileInfo{
			FName:    "test",
			FModTime: time.Now(),
			FIsDir:   false,
			FMode:    os.ModeType,
		},
		FileObject: model.FileObject{
			Bucket: "test",
			Key:    "test",
		},
	}
	testFile1WithNewContent := &model.FileDetail{
		FileInfo: model.FileInfo{
			FName:    "test",
			FModTime: time.Now(),
			FIsDir:   false,
			FMode:    os.ModeType,
		},
		FileObject: model.FileObject{
			Bucket: "newcontent",
			Key:    "newcontent",
		},
	}
	testFile2 := &model.FileDetail{
		FileInfo: model.FileInfo{
			FName:    "bbb",
			FModTime: time.Now(),
			FIsDir:   false,
			FMode:    os.ModeType,
		},
		FileObject: model.FileObject{
			Bucket: "bbb",
			Key:    "bbb",
		},
	}

	err = txn.AddFile(ctx, userId, "/test", testFile1)
	assert.NoError(err)

	err = txn.AddFile(ctx, userId, "/test/test", testFile1)
	assert.Error(err)

	err = txn.AddFile(ctx, userId, "/test", testFile1WithNewContent)
	assert.NoError(err)

	err = txn.AddFile(ctx, userId, "/bbb", testFile2)
	assert.NoError(err)

	var info *model.FileInfoWrapper
	info, err = txn.GetFileInfo(ctx, userId, "/test")
	assert.NoError(err)

	assert.Equal(&model.FileInfoWrapper{
		FileInfo:   testFile1WithNewContent.FileInfo,
		FileObject: &testFile1WithNewContent.FileObject,
	}, info)

	info, err = txn.GetFileInfo(ctx, userId, "/")
	assert.NoError(err)

	assert.ElementsMatch(model.DirContentList{
		*testFile1WithNewContent.AsDirContentItem(), *testFile2.AsDirContentItem(),
	}, info.DirContent)

}
