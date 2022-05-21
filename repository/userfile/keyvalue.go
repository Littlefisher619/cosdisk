package userfile

import (
	"os"
	"path"
	"strings"

	"github.com/Littlefisher619/cosdisk/model"
)

type FileData = model.FileData

/* help functions */

func GenerateUserFileKey(username string, path string) string {
	return username + ":" + path
}

func ParseUserFileKey(key string) (username string, path string, err error) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		err = model.ErrInvalidFileKey
		return
	}
	username = parts[0]
	path = parts[1]
	return
}

const (
	FileType int = 0
	DirType  int = 1
)

func GenerateUserFileValue(id string) string {
	return "file" + ":" + id
}

func GenerateUserDirValue(files []string) string {
	return "dir" + ":" + strings.Join(files, ";")
}

func ParseUserFileValue(value string) (kind int, data []string, err error) {
	parts := strings.Split(value, ":")
	err = nil
	if len(parts) != 2 {
		err = model.ErrInvalidFileKey
		return
	}
	switch parts[0] {
	case "file":
		kind = FileType
		data = append(data, parts[1])
	case "dir":
		kind = DirType
		if parts[1] == "" {
			data = nil
			return
		}
		data = append(data, strings.Split(parts[1], ";")...)
	}
	return
}

func ParseDirContent(data []string) (files []os.FileInfo, err error) {
	for _, file := range data {
		if file[len(file)-1] == '/' {
			files = append(files, &FileData{
				Minzhi:     file[:len(file)-1],
				Wenjianjia: true,
				Mushi:      os.ModeDir,
			})
		} else {
			files = append(files, &FileData{
				Minzhi:     file,
				Wenjianjia: false,
			})
		}
	}
	return files, nil
}

func Makefilepath(filepath string) (parentpath string, filename string) {
	filename = path.Base(filepath)
	parentpath = path.Dir(filepath)
	if parentpath != "/" {
		parentpath = parentpath + "/"
	}
	return
}

func Makedirpath(dirpath string) (parentpath string, dirname string) {
	dirname = path.Base(dirpath) + "/"
	if dirpath[len(dirpath)-1] == '/' {
		parentpath = path.Dir(dirpath[:len(dirpath)-1])
	} else {
		parentpath = path.Dir(dirpath)
	}
	if parentpath != "/" {
		parentpath = parentpath + "/"
	}
	return
}
