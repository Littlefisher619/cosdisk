package ftp

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// File is the afero.File implementation
type File struct {
	Fs       *FileSystem
	path     string
	fileData os.FileInfo
	content  []byte
	At       int64
	isWrited bool
	isCached bool
}

// Name of the file
func (f *File) Name() string { return f.fileData.Name() }

// Close closes the file transfer and does the mail sending
func (f *File) Close() error {
	f.Fs.logger.Info("Close " + f.path)
	if f.Fs == nil {
		return ErrNotFound
	}
	if f.isWrited {
		f.At = 0
		err := f.Fs.cosdisk.UploadUserFileByReader(context.Background(), f.Fs.userId, f.path, f)
		if err != nil {
			return err
		}
	}
	f.isCached = false
	f.content = []byte{}
	f.At = 0
	delete(f.Fs.openfiles, f.path)
	return nil
}

// Read stores the received file content into the local buffer
func (f *File) Read(b []byte) (int, error) {
	f.Fs.logger.Info("Read " + f.path)
	if !f.isCached {
		reader, err := f.Fs.cosdisk.DownloadUserFileByReader(context.Background(), f.Fs.userId, f.path)
		if err != nil {
			return 0, err
		}
		f.content, err = ioutil.ReadAll(reader)
		if err != nil {
			return 0, err
		}
		f.isCached = true
		f.isWrited = false
		f.At = 0
	}
	n := 0
	if len(b) > 0 && int(f.At) == len(f.content) {
		return 0, io.EOF
	}
	if len(f.content)-int(f.At) >= len(b) {
		n = len(b)
	} else {
		n = len(f.content) - int(f.At)
	}
	copy(b, f.content[f.At:f.At+int64(n)])
	f.At += int64(n)

	return n, nil
}

// ReadAt is not implemented
func (f *File) ReadAt(_ []byte, _ int64) (int, error) {
	f.Fs.logger.Info("ReadAt " + f.path)
	return 0, ErrNotImplemented
}

// Truncate is not implemented
func (f *File) Truncate(size int64) error {
	f.Fs.logger.Info("Truncate " + f.path)
	return ErrNotImplemented
}

func (f *File) Readdir(count int) (result []os.FileInfo, err error) {
	f.Fs.logger.Info("Readdir " + f.path)

	res, err := f.Fs.cosdisk.ListFiles(f.Fs.userId, f.path)
	if err != nil {
		f.Fs.logger.Error("readdir failed: " + err.Error())
	}
	return res, nil
}

// Readdirnames is not implemented
func (f *File) Readdirnames(_ int) ([]string, error) {
	f.Fs.logger.Info("Readdirnames " + f.path)
	return []string{}, ErrNotImplemented
}

// Seek is not implemented
func (f *File) Seek(offset int64, whence int) (ret int64, err error) {
	f.Fs.logger.Info("Seek " + f.path + " " + fmt.Sprint(offset) + " " + fmt.Sprint(whence))

	if whence == 0 {
		f.At = offset
	} else if whence == 1 {
		f.At += offset
	} else if whence == 2 {
		f.At = int64(len(f.content)) + offset
	}
	return 0, nil
}

// Stat is not implemented
func (f *File) Stat() (os.FileInfo, error) {
	return f.fileData, nil
}

// Sync is not implemented
func (f *File) Sync() error {
	return nil
}

// WriteString is not implemented
func (f *File) WriteString(s string) (int, error) {
	return 0, ErrNotImplemented
}

// WriteAt is not implemented
func (f *File) WriteAt(b []byte, off int64) (int, error) {
	return 0, ErrNotImplemented
}

func (f *File) Write(b []byte) (int, error) {
	f.Fs.logger.Info("write " + f.path + " " + string(b))
	f.content = append(f.content[:f.At], b...)
	f.At += int64(len(b))
	f.isWrited = true
	f.isCached = true

	return len(b), nil
}
