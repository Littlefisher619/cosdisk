package model

import (
	"os"
	"time"
)

func (fd FileInfo) Sys() interface{} {
	return nil
}

func (fd FileInfo) Name() string {
	return fd.FName
}

func (fd FileInfo) Size() int64 {
	return fd.FSize
}

func (fd FileInfo) IsDir() bool {
	return fd.FIsDir
}

func (fd FileInfo) Mode() os.FileMode {
	if fd.FIsDir {
		return os.ModeDir
	}
	return os.ModeType
}

// modification time
func (fd FileInfo) ModTime() time.Time {
	return fd.FModTime
}

func (fdl DirContentList) AsOsFileInfo() []os.FileInfo {
	var fis []os.FileInfo = make([]os.FileInfo, len(fdl))
	for i := range fdl {
		fis[i] = fdl[i]
	}
	return fis
}
