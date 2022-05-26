package dbdriver

import (
	"context"
)

// KeyValueStorage is the interface for kv storage
// Eg redis, tikv, etc
// The storage engine should implement this interface, and start a KeyValueTXN
type KVTxn interface {
	Commit(ctx context.Context) (err error)
	Rollback() (err error)
	KVOperation
}

type KVOperation interface {
	Set(key, value Serializable) (err error)
	Get(ctx context.Context, key Serializable, value Serializable) (err error)
	Delete(key Serializable) (err error)
}

type KVStorage interface {
	Begin() (KVTxn, error)
	Close() error
}

type Serializable interface {
	Serialize() ([]byte, error)
	UnSerialize([]byte) error
}
