package userfile

import (
	"context"
	"fmt"
	kv "github.com/Littlefisher619/cosdisk/pkg/metadatakv"
)

type mutation struct {
	existV *kv.Value

	parentK *kv.Key
	parentV *kv.Value
}

type validator interface {
	validate(v *kv.Value) error
}

type existOption byte
type fileTypeOption byte

type mustEmpty struct{}

func (x mustEmpty) validate(v *kv.Value) error {
	if len(v.DirTypeValue.DirContentList) != 0 {
		return ErrDirectoryNotEmpty
	}
	return nil
}

const (
	mustNil = existOption(iota)
	mustNotNil
)

const (
	mustFileOrNil = fileTypeOption(iota)
	mustDirOrNil
)

func (x existOption) validate(v *kv.Value) error {
	if x == mustNotNil {
		if v == nil {
			return ErrNotFound
		}
		return nil
	}
	if x == mustNil {
		if v == nil {
			return ErrPathExists
		}
		return nil
	}

	return nil
}

func (x fileTypeOption) validate(v *kv.Value) error {
	if v == nil {
		return nil
	}
	if x == mustDirOrNil {
		if !v.IsDir() {
			return ErrIsFile
		}
		return nil
	}

	if x == mustFileOrNil {
		if v.IsDir() {
			return ErrIsDir
		}
		return nil
	}

	return nil
}

// buildMutation will check existence of inputKey by options
// and get information of its parent value
// when parent not found or is a file, ErrParentDirNotExist or ErrIsFile will returned
func (txn *Txn) buildMutation(
	ctx context.Context,
	inputK *kv.Key,
	optionsOnExistV ...validator,
) (*mutation, error) {
	existV, err := txn.getByKey(ctx, inputK)
	if err != nil {
		return nil, fmt.Errorf("fail to check key exist: %w", err)
	}

	for _, o := range optionsOnExistV {
		err = o.validate(existV)
		if err != nil {
			return nil, err
		}
	}

	parentK, parentV, err := txn.getParentOfKey(ctx, inputK)
	if err != nil {
		return nil, fmt.Errorf("fail to get parent dir: %w", err)
	}
	if parentV == nil {
		return nil, ErrParentDirNotExist
	}
	if !parentV.IsDir() {
		if existV != nil {
			return nil, fmt.Errorf("Unexcepted error: file exist on a file")
		}
		return nil, ErrIsFile
	}

	return &mutation{
		existV:  existV,
		parentK: parentK,
		parentV: parentV,
	}, nil
}
