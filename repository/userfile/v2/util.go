package userfile

import (
	"context"
	"fmt"
	"path"

	"github.com/Littlefisher619/cosdisk/pkg/dbdriver"
	kv "github.com/Littlefisher619/cosdisk/pkg/metadatakv"
)

func makeKey(UID, inputPath string) (
	*kv.Key, error,
) {
	k := &kv.Key{
		UID:  UID,
		Path: inputPath,
	}
	if err := k.Cleanup(); err != nil {
		return nil, malformedInputError(err)

	}
	return k, nil
}

func (txn *Txn) getByKey(ctx context.Context, k *kv.Key) (*kv.Value, error) {
	value := &kv.Value{}
	err := txn.Get(ctx, k, value)
	if err != nil {
		if err == dbdriver.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("DB Error: %w", err)
	}
	return value, nil
}

func (txn *Txn) getFromInput(ctx context.Context, UID, inputPath string) (*kv.Value, error) {
	k, err := makeKey(UID, inputPath)
	if err != nil {
		return nil, err
	}
	return txn.getByKey(ctx, k)
}

func (txn *Txn) getParentOfKey(ctx context.Context, k *kv.Key) (
	*kv.Key, *kv.Value, error,
) {
	parentPath := path.Dir(k.Path)
	parentK := &kv.Key{UID: k.UID, Path: parentPath}

	parentV, err := txn.getByKey(ctx, parentK)
	if err != nil {
		return parentK, nil, fmt.Errorf("get parent dir: %w", err)
	}

	if parentV == nil {
		return parentK, nil, nil
	}

	return parentK, parentV, nil
}

func malformedInputError(e error) error {
	return fmt.Errorf("malformed input: %s", e)
}
