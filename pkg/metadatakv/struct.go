package kv

import (
	"github.com/Littlefisher619/cosdisk/model"
)

type ValueType byte
type FileValue = model.FileDetail
type DirValue = model.DirDetail

const (
	keySperator = byte(0x00)

	FileValueType = ValueType(0x01)
	DirValueType  = ValueType(0x00)
)

type Key struct {
	UID, Path string
}

type Value struct {
	Type ValueType

	FileTypeValue *FileValue
	DirTypeValue  *DirValue
}


