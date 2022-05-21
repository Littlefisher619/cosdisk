package userfile

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"

	model "github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/repository/connections"
)

func TestAll(t *testing.T) {
	var userfile model.UserfileRepository
	rc := connections.InitTikv([]string{"127.0.0.1:2379"})
	if rc == nil {
		return
	}
	defer rc.CloseTikv()
	userfile = NewUserfileKV(rc)
	//userfile = NewMap()
	txn, err := userfile.StartTranscation()
	if err != nil {
		t.Fatalf("StartTranscation failed")
	}

	t.Run("TestAddandListFile", func(t *testing.T) {
		ass := require.New(t)

		userid := "stupidfish"
		// test add not root file
		ass.Error(txn.AddFile(userid, "wrong", "456"))

		txn.DeleteDir(userid, "/file2.txt/")
		txn.DeleteDir(userid, "/file2.txt")
		txn.DeleteDir(userid, "/test1/")
		txn.DeleteFile(userid, "/test1")
		txn.DeleteDir(userid, "/test2/")
		txn.DeleteDir(userid, "/")
		txn.DeleteDir(userid, "/test2/test23")
		txn.DeleteFile(userid, "/test1/file2.txt")

		// all add file and dir must under "/"
		ass.NoError(txn.AddDir(userid, "/"))
		ass.NoError(txn.AddDir(userid, "/test1/"))
		ass.NoError(txn.AddFile(userid, "/test1/file1.txt", "123"))
		ass.NoError(txn.AddFile(userid, "/file2.txt", "456"))
		// update file
		ass.NoError(txn.AddFile(userid, "/file2.txt", "456"))
		ass.NoError(txn.AddFile(userid, "/file2.txt", "789"))

		ass.NoError(txn.AddDir(userid, "/test2"))
		ass.NoError(txn.AddDir(userid, "/test2/test23"))

		// test add dir again
		ass.Error(txn.AddDir(userid, "/"))
		ass.Error(txn.AddDir(userid, "/test1/"))
		ass.Error(txn.AddDir(userid, "/test1"))

		// add same file name as dir
		ass.Error(txn.AddDir(userid, "/test1/file1.txt"))

		// add same dir name as file
		ass.Error(txn.AddFile(userid, "/test1", "123"))

		// add at wrong path
		ass.Error(txn.AddFile(userid, "/file2.txt/filex.txt", "123"))

		files, err := txn.ListFiles(userid, "/")
		ass.NoError(err)
		dir, err := ParseDirContent([]string{"test1/", "file2.txt", "test2/"})
		ass.NoError(err)
		ass.Equal(dir, files)

		files, err = txn.ListFiles(userid, "/test1/")
		ass.NoError(err)
		dir, err = ParseDirContent([]string{"file1.txt"})
		ass.NoError(err)
		ass.Equal(dir, files)
		ass.EqualValues(dir, []fs.FileInfo{&FileData{
			Minzhi:     "file1.txt",
			Wenjianjia: false,
		}})

		files, err = txn.ListFiles(userid, "/test1")
		ass.NoError(err)
		ass.Equal(dir, files)

		_, err = txn.ListFiles(userid, "/test2/")
		ass.NoError(err)

		_, err = txn.ListFiles(userid, "/test3")
		ass.Equal(err, model.ErrFileNotFound)

		id, err := txn.GetFileID(userid, "/test1/file1.txt")
		ass.NoError(err)
		ass.Equal("123", id)

		id, err = txn.GetFileID(userid, "/file2.txt")
		ass.NoError(err)
		ass.Equal("789", id)
		//rc.CloseTikv()
	})

	t.Run("TestGetInfo", func(t *testing.T) {
		ass := require.New(t)

		userid := "stupidfish"
		info, err := txn.GetFileInfo(userid, "/")
		ass.NoError(err)
		ass.Equal(info.Name(), "/")
		ass.Equal(info.IsDir(), true)

		info, err = txn.GetFileInfo(userid, "/test1")
		ass.NoError(err)
		ass.Equal(info.Name(), "test1")
		ass.Equal(info.IsDir(), true)
		_, err = txn.GetFileInfo(userid, "/test1/")
		ass.NoError(err)

		info, err = txn.GetFileInfo(userid, "/test2")
		ass.NoError(err)
		ass.Equal(info.Name(), "test2")
		ass.Equal(info.IsDir(), true)
		_, err = txn.GetFileInfo(userid, "/test2/")
		ass.NoError(err)

		info, err = txn.GetFileInfo(userid, "/file2.txt")
		ass.NoError(err)
		ass.Equal(info.Name(), "file2.txt")
		ass.Equal(info.IsDir(), false)

		_, err = txn.GetFileInfo(userid, "/file3.txt")
		ass.Equal(err, model.ErrFileNotFound)
	})

	t.Run("TestDeleteFile", func(t *testing.T) {
		ass := require.New(t)

		userid := "stupidfish"
		ass.NoError(txn.DeleteFile(userid, "/test1/file1.txt")) // remove file

		files, err := txn.ListFiles(userid, "/test1/")
		ass.NoError(err)
		ass.Equal([]fs.FileInfo(nil), files)

		ass.NoError(txn.AddFile(userid, "/test1/file2.txt", "123"))

		_, err = txn.GetFileID(userid, "/test1/file1.txt")
		ass.Error(err)

		files, err = txn.ListFiles(userid, "/test1/")
		ass.NoError(err)
		dir, err := ParseDirContent([]string{"file2.txt"})
		ass.NoError(err)
		ass.Equal(dir, files)

		// test delete wrong type
		ass.Error(txn.DeleteFile(userid, "/test1"))

		ass.NoError(txn.DeleteFile(userid, "/test1/file2.txt"))
		ass.NoError(txn.DeleteDir(userid, "/test1")) // remove dir
		ass.NoError(txn.DeleteDir(userid, "/test2/test23"))
		ass.NoError(txn.DeleteDir(userid, "/test2/")) // remove dir

		// delete dir again
		ass.Equal(txn.DeleteDir(userid, "/test1"), model.ErrFileNotFound)
		// delete file again
		ass.Error(txn.DeleteFile(userid, "/test1/file1.txt")) // remove file
		// test delete wrong type
		ass.Error(txn.DeleteDir(userid, "/file2.txt")) // remove dir
		// test delete wrong path
		ass.Error(txn.DeleteDir(userid, "/file2.txt/dirx")) // remove dir

		ass.NoError(txn.DeleteFile(userid, "/file2.txt")) // remove file

		files, err = txn.ListFiles(userid, "/")
		ass.NoError(err)
		ass.Equal([]fs.FileInfo(nil), files)

		ass.NoError(txn.DeleteDir(userid, "/")) // remove root
	})
	err = txn.CommitTranscation()
	if err != nil {
		t.Fatalf("CommitTranscation failed")
	}
}
