package ftp

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Littlefisher619/cosdisk/model"
	cosservice "github.com/Littlefisher619/cosdisk/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// ErrNotImplemented is returned when something is not implemented
var ErrNotImplemented = errors.New("not implemented")

// ErrNotFound is returned when something is not found
var ErrNotFound = errors.New("not found")

// ErrInvalidParameter is returned when a parameter is invalid
var ErrInvalidParameter = errors.New("invalid parameter")

// The ClientDriver is the internal structure used for handling the client. At this stage it's limited to the afero.Fs
type ClientDriver struct {
	afero.Fs
}

type FileSystem struct {
	userId    string
	cwd       string
	openfiles map[string]*File
	cosdisk   *cosservice.CosDisk
	logger    *logrus.Entry
}

// Name of the filesystem
func (f *FileSystem) Name() string {
	return "cosdisk-filesystem"
}

// Chtimes is not implemented
func (f *FileSystem) Chtimes(name string, atime, mtime time.Time) error {
	return ErrNotImplemented
}

// Chmod is not implemented
func (f *FileSystem) Chmod(name string, mode os.FileMode) error {
	f.logger.Info("Chmod " + name)
	return ErrNotImplemented
}

// Rename is not implemented
func (f *FileSystem) Rename(name string, newname string) error {
	f.logger.Info("Rename " + name + " to " + newname)
	return ErrNotImplemented
}

// Chown is not implemented
func (f *FileSystem) Chown(string, int, int) error {
	return ErrNotImplemented
}

// RemoveAll is not implemented
func (f *FileSystem) RemoveAll(name string) error {
	f.logger.Info("RemoveAll " + name)
	return ErrNotImplemented
}

func (f *FileSystem) Remove(name string) error {
	f.logger.Info("Remove " + name)

	data, err := f.Stat(name)
	if err != nil {
		return err
	}
	if data.IsDir() {
		return f.cosdisk.DeleteDir(f.userId, name)
	} else {
		return f.cosdisk.DeleteFIle(f.userId, name)
	}
}

// Mkdir is not implemented
func (f *FileSystem) Mkdir(name string, mode os.FileMode) error {
	f.logger.Info("Mkdir " + name)

	return f.cosdisk.CreateDir(f.userId, name)
	//return os.Mkdir(name, mode)
}

// MkdirAll is not implemented
func (f *FileSystem) MkdirAll(name string, mode os.FileMode) error {
	f.logger.Info("MkdirAll " + name)
	return ErrNotImplemented
}

// Open opens a file buffer
func (f *FileSystem) Open(name string) (afero.File, error) {
	f.logger.Info("Open " + name)

	name = strings.TrimSuffix(name, "-a")
	if file, ok := f.openfiles[name]; ok {
		return file, nil
	}
	data, err := f.Stat(name)
	if err != nil {
		return nil, err
	}
	file := &File{fileData: data, Fs: f, path: name}
	f.openfiles[name] = file
	return file, nil

	//return os.Open(name)
}

// Create creates a file buffer
func (f *FileSystem) Create(name string) (afero.File, error) {
	f.logger.Info("Create " + name)

	file := &File{fileData: &model.FileData{Minzhi: path.Base(name), Wenjianjia: false},
		Fs: f, path: name, isWrited: true, isCached: true}
	f.openfiles[name] = file
	return file, nil
}

// OpenFile opens a file buffer
func (f *FileSystem) OpenFile(name string, flag int, mode os.FileMode) (afero.File, error) {
	f.logger.Info("OpenFile "+name, flag, mode)
	file, err := f.Open(name)
	if err != nil {
		if err == model.ErrFileNotFound {
			if flag&os.O_CREATE != 0 {
				return f.Create(name)
			}
			return nil, os.ErrNotExist
		}
	}
	return file, nil
}

func (f *FileSystem) Stat(name string) (os.FileInfo, error) {
	f.logger.Info("Stat " + name)
	/*
		info, err := os.Stat(name)
		f.logger.Info(info)
		f.logger.Info(err)
		return info, err
	*/
	if strings.HasSuffix(name, "-a") {
		return nil, &os.PathError{Path: name, Err: os.ErrNotExist}
	}
	if file, ok := f.openfiles[name]; ok {
		return file.fileData, nil
	} else {
		f.cwd = name
		res, err := f.cosdisk.GetFileInfo(f.userId, name)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

}

// LstatIfPossible is not implemented
func (f *FileSystem) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	f.logger.Info("LstatIfPossible " + name)
	return nil, false, &os.PathError{Op: "lstat", Path: name, Err: nil}
}
