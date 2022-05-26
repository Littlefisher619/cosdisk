package userfile

import (
	"context"

	"github.com/Littlefisher619/cosdisk/model"
)

func (txn *Txn) ListFiles(ctx context.Context, UID, inputPath string) (files model.DirContentList, err error) {
	v, err := txn.getFromInput(ctx, UID, inputPath)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNotFound
	}
	if !v.IsDir() {
		return nil, ErrIsFile
	}

	return v.DirTypeValue.DirContentList, nil
}

func (txn *Txn) GetFileInfo(ctx context.Context, UID, inputPath string) (f *model.FileInfoWrapper, err error) {
	k, err := makeKey(UID, inputPath)
	if err != nil {
		return nil, malformedInputError(err)
	}
	v, err := txn.getByKey(ctx, k)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNotFound
	}

	if v.IsDir() {
		return &model.FileInfoWrapper{
			FileInfo:   v.DirTypeValue.FileInfo,
			DirContent: v.DirTypeValue.DirContentList,
		}, nil
	}
	return &model.FileInfoWrapper{
		FileInfo:   v.FileTypeValue.FileInfo,
		FileObject: &v.FileTypeValue.FileObject,
	}, nil
}
