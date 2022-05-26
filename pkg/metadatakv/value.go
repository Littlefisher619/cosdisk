package kv

import (
	"os"
	"time"
)

func (fd *Value) Name() string {
	if fd.IsDir() {
		return fd.DirTypeValue.FileInfo.FName
	} else {
		return fd.FileTypeValue.FileInfo.FName
	}
}

func (fd *Value) Sys() interface{} {
	return nil
}

func (fd *Value) IsDir() bool {
	return fd.Type == DirValueType
}
func (fd *Value) Size() int64 {
	if fd.IsDir() {
		return fd.DirTypeValue.FileInfo.FSize
	} else {
		return fd.FileTypeValue.FileInfo.FSize
	}
}

func (fd *Value) Mode() os.FileMode {
	if fd.IsDir() {
		return os.ModeDir
	} else {
		return os.ModeType
	}
}

// modification time
func (fd *Value) ModTime() time.Time {
	if fd.IsDir() {
		return fd.DirTypeValue.FileInfo.FModTime
	} else {
		return fd.FileTypeValue.FileInfo.FModTime
	}
}
