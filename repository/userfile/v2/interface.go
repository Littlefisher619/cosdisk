package userfile

// import (
// 	"context"
// )

// KeyValueStorage is the interface for kv storage
// Eg redis, tikv, etc
// The storage engine should implement this interface, and start a KeyValueTXN
// type KeyValueStorage interface {
// 	StartTranscation() (KeyValueTXN, error)
// }

// // KeyValueTXN is a transaction for key-value storage
// type KeyValueTXN interface {
// 	Set(key, value Serializable) (err error)
// 	Get(ctx context.Context, key Serializable, value Serializable) (err error)
// 	Delete(key Serializable) (err error)
// 	CommitTranscation() (err error)
// 	RollingBackTranscation() (err error)
// }

// type Serializable interface {
// 	Serialize() ([]byte, error)
// 	UnSerialize([]byte) error
// }