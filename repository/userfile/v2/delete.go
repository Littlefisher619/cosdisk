package userfile

import (
	"context"
	"time"

	kv "github.com/Littlefisher619/cosdisk/pkg/metadatakv"
)

func (txn *Txn) removeSummaryFromDirValue(
	targetSummaryFileName string,
	dirK *kv.Key, dirV *kv.Value,
) error {

	var index = -1
	for i, v := range dirV.DirTypeValue.DirContentList {
		if v.FName == targetSummaryFileName {
			index = i
			break
		}
	}
	if index == -1 {
		return ErrNotFound
	}
	dirV.DirTypeValue.DirContentList = append(dirV.DirTypeValue.DirContentList[:index], dirV.DirTypeValue.DirContentList[index+1:]...)
	dirV.DirTypeValue.FileInfo.FModTime = time.Now()
	return txn.Set(dirK, dirV)

}

func (txn *Txn) DeleteDir(ctx context.Context, UID, dirPathToRemove string) error {
	dirToRemoveK, err := makeKey(UID, dirPathToRemove)
	if err != nil {
		return malformedInputError(err)
	}

	if dirToRemoveK.Path == "/" {
		return ErrDeleteRootDir
	}

	op, err := txn.buildMutation(ctx, dirToRemoveK, mustDirOrNil, mustNotNil)
	if err != nil {
		return err
	}

	err = txn.removeSummaryFromDirValue(op.existV.Name(), op.parentK, op.parentV)
	if err != nil {
		return err
	}

	err = txn.KVTxn.Delete(dirToRemoveK)
	if err != nil {
		return err
	}

	return nil
}

func (txn *Txn) DeleteFile(ctx context.Context, UID, inputPath string) error {
	fileToRemoveK, err := makeKey(UID, inputPath)
	if err != nil {
		return err
	}

	op, err := txn.buildMutation(ctx, fileToRemoveK, mustFileOrNil, mustNotNil)
	if err != nil {
		return err
	}

	err = txn.removeSummaryFromDirValue(op.existV.Name(), op.parentK, op.parentV)
	if err != nil {
		return err
	}

	err = txn.Delete(fileToRemoveK)
	if err != nil {
		return err
	}

	return nil
}
