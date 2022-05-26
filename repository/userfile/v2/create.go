package userfile

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Littlefisher619/cosdisk/model"
	kv "github.com/Littlefisher619/cosdisk/pkg/metadatakv"
)

func (txn *Txn) CreateRepository(ctx context.Context, UID string) (
	err error,
) {
	rootK := &kv.Key{
		UID:  UID,
		Path: "/",
	}

	rootV := &kv.Value{
		Type: kv.DirValueType,
		DirTypeValue: &kv.DirValue{
			FileInfo: model.FileInfo{
				FName:    "/",
				FModTime: time.Now(),
				FIsDir:   true,
				FMode:    os.ModeDir,
			},
		},
	}

	return txn.Set(rootK, rootV)
}

func (txn *Txn) AddFile(
	ctx context.Context, UID, inputPath string,
	file *model.FileDetail,
) error {
	inputK, err := makeKey(UID, inputPath)
	if err != nil {
		return malformedInputError(err)
	}

	op, err := txn.buildMutation(ctx, inputK, mustFileOrNil)
	if err != nil {
		return err
	}

	if op.existV != nil {
		// overwrite file with file, permits!
		op.existV.FileTypeValue.FileObject = file.FileObject
		op.existV.FileTypeValue.FileInfo = file.FileInfo
		return txn.KVTxn.Set(inputK, op.existV)
	}

	// Create
	newFileV := &kv.Value{
		Type:          kv.FileValueType,
		FileTypeValue: file,
	}
	err = txn.KVTxn.Set(inputK, newFileV)

	if err != nil {
		return fmt.Errorf("fail to create file record: %w", err)
	}

	addItemToParentDirList(newFileV.FileTypeValue.AsDirContentItem(), op.parentV)

	err = txn.KVTxn.Set(op.parentK, op.parentV)

	if err != nil {
		return fmt.Errorf("fail to modify parent dir record: %w", err)
	}

	return nil
}

func (txn *Txn) AddDir(ctx context.Context, UID, destPath string) error {
	inputK, err := makeKey(UID, destPath)
	if err != nil {
		return malformedInputError(err)
	}

	op, err := txn.buildMutation(
		ctx, inputK, mustNil,
	)

	// Create
	newDirV := &kv.Value{
		Type: kv.DirValueType,
		DirTypeValue: &kv.DirValue{
			FileInfo: model.FileInfo{
				FModTime: time.Now(),
				FName:    path.Base(inputK.Path),
				FIsDir:   true,
				FMode:    os.ModeDir,
			},
		},
	}
	err = txn.KVTxn.Set(inputK, newDirV)

	if err != nil {
		return fmt.Errorf("fail to create dir record: %w", err)
	}

	addItemToParentDirList(newDirV.DirTypeValue.AsDirContentItem(), op.parentV)

	err = txn.KVTxn.Set(op.parentK, op.parentV)

	if err != nil {
		return fmt.Errorf("fail to modify parent dir record: %w", err)
	}

	return nil
}

func addItemToParentDirList(
	valToAdd *model.DirContentItem, parentDirV *kv.Value,
) {
	parentDirV.DirTypeValue.FileInfo.FModTime = time.Now()
	parentDirV.DirTypeValue.DirContentList = append(parentDirV.DirTypeValue.DirContentList, *valToAdd)
}
