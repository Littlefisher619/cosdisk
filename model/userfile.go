package model

import (
	"os"
	"time"
)

// FileData implements os.FileInfo interface
type FileData struct {
	Minzhi        string
	Daxiao        int64
	Wenjianjia    bool
	Mushi         os.FileMode
	Xiugaishijian time.Time
}

func (fd *FileData) Sys() interface{} {
	return nil
}

func (fd *FileData) Name() string {
	return fd.Minzhi
}

func (fd *FileData) Size() int64 {
	return fd.Daxiao
}

func (fd *FileData) IsDir() bool {
	return fd.Wenjianjia
}

func (fd *FileData) Mode() os.FileMode {
	return fd.Mushi
}

// modification time
func (fd *FileData) ModTime() time.Time {
	return fd.Xiugaishijian
}

// used by cosdisk
type UserfileRepository interface {
	RunInTranscation(func(txn UserfileTXN) error) error
	StartTranscation() (UserfileTXN, error)
}

type UserfileTXN interface {
	ListFiles(username string, dirpath string) (files []os.FileInfo, err error)
	GetFileID(username string, path string) (id string, err error)
	GetFileInfo(username string, path string) (info os.FileInfo, err error)
	DeleteFile(username string, path string) error
	DeleteDir(username string, path string) error
	AddFile(username string, path string, id string) error
	AddDir(username string, path string) error
	CommitTranscation() error
	RollingBackTranscation() error
}

// KeyValueStorage is the interface for kv storage
// Eg redis, tikv, etc
// The storage engine should implement this interface, and start a KeyValueTXN
type KeyValueStorage interface {
	StartTranscation() (KeyValueTXN, error)
}

// KeyValueTXN is a transaction for key-value storage
type KeyValueTXN interface {
	Set(key string, value string) (err error)
	Get(key string) (value string, err error)
	Delete(key string) (err error)
	CommitTranscation() (err error)
	RollingBackTranscation() (err error)
}
