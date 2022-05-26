package kv

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

func (v *Value) Serialize() ([]byte, error) {
	return msgpack.Marshal(v)
	var (
		b   []byte
		err error
	)
	switch v.Type {
	case FileValueType:
		b, err = msgpack.Marshal(v.FileTypeValue)
	case DirValueType:
		b, err = msgpack.Marshal(v.DirTypeValue)
	default:
		return nil, fmt.Errorf("invalid value type: got %d, wanted 0x01 or 0x00", v.Type)
	}
	if err != nil {
		return nil, err
	}

	return append([]byte{byte(v.Type)}, b...), nil
}

func (v *Value) UnSerialize(b []byte) error {
	newV := Value{}
	err := msgpack.Unmarshal(b, &newV)
	if err != nil {
		return err
	}
	*v = newV
	return nil

	if len(b) < 1 {
		return fmt.Errorf("invalid data length: got %d, wanted >=1", len(b))
	}
	switch ValueType(b[0]) {
	case FileValueType:
		v.Type = FileValueType
		v.FileTypeValue = &FileValue{}
		return msgpack.Unmarshal(b[1:], v.FileTypeValue)
	case DirValueType:
		v.Type = DirValueType
		v.DirTypeValue = &DirValue{}
		return msgpack.Unmarshal(b[1:], v.DirTypeValue)
	default:
		return fmt.Errorf("invalid value type: got %d, wanted 0x01 or 0x00", b[0])
	}
}

func (k *Key) Cleanup() error {
	if strings.IndexByte(k.Path, keySperator) != -1 {
		return fmt.Errorf("path contains invalid character")
	}
	if !strings.HasPrefix(k.Path, "/") {
		return fmt.Errorf("path should be absolute")
	}
	if strings.IndexByte(k.UID, keySperator) != -1 {
		return fmt.Errorf("username contains invalid character")
	}
	k.Path = path.Clean(k.Path)
	return nil
}

func (v *Key) Serialize() ([]byte, error) {
	b := make([]byte, len(v.UID)+len(v.Path)+1)
	copy(b, v.UID)
	copy(b[len(v.UID):], []byte{keySperator})
	copy(b[len(v.UID)+1:], v.Path)
	return b, nil
}

func (v *Key) UnSerialize(b []byte) error {
	idx := bytes.IndexByte(b, keySperator)
	if idx == -1 {
		return fmt.Errorf("invalid data: no separator found")
	}
	v.UID = string(b[:idx])
	v.Path = string(b[idx+1:])
	return nil
}
