package model

import (
	"context"
	"os"
	"time"
)

// FileInfo implements os.FileInfo interface
type FileInfo struct {
	FSize    int64
	FMode    os.FileMode
	FModTime time.Time
	FIsDir   bool
	FName    string
}

type TxnFunc func(txn UserfileTXN) error

// used by cosdisk
type UserfileRepository interface {
	RunInTxn(ctx context.Context, txnfunc TxnFunc) error
	Begin() (UserfileTXN, error)
	Close() error
}

type DirContentList []DirContentItem

type UserfileOperation interface {
	CreateRepository(ctx context.Context, UID string) error
	ListFiles(ctx context.Context, UID, inputPath string) (DirContentList, error)
	// ListFileSummaries(ctx context.Context, UID, inputPath string) (FileDataList, error)
	GetFileInfo(ctx context.Context, UID, inputPath string) (*FileInfoWrapper, error)
	DeleteFile(ctx context.Context, UID, inputPath string) error
	DeleteDir(ctx context.Context, UID, inputPath string) error
	// AddFile takes a filepath and a file data
	//
	AddFile(
		ctx context.Context, UID, inputPath string,
		file *FileDetail,
	) error
	AddDir(ctx context.Context, UID, inputPath string) error
}
type FileInfoWrapper struct {
	FileInfo
	// Only valid when is FileInfo.IsDir() == true
	DirContent DirContentList
	// Only valid when is FileInfo.IsDir() == false
	FileObject *FileObject
}

type DirContentItem struct {
	FileInfo
	// Only valid when is FileDetail.IsDir() == false
	*FileObject
}

type DirDetail struct {
	FileInfo
	DirContentList
}

func (d* DirDetail) AsDirContentItem() *DirContentItem {
	return &DirContentItem{
		FileInfo: d.FileInfo,
		FileObject: nil,
	}
}

type FileDetail struct {
	FileInfo
	FileObject
}

func (d* FileDetail) AsDirContentItem() *DirContentItem {
	return &DirContentItem{
		FileInfo: d.FileInfo,
		FileObject: &d.FileObject,
	}
}


type UserfileTXN interface {
	UserfileOperation
	Commit(ctx context.Context) error
	Rollback() error
}

type FileObject struct {
	Bucket string
	Key    string
}
