package dbdriver

import (
	"context"
	"fmt"
)

type MapKVClient map[string][]byte

type MapKVTxn struct {
	MapKVClient
}

func InitMapKV() *MapKVClient {
	return &MapKVClient{}
}

func (txn *MapKVTxn) Set(key Serializable, value Serializable) (err error) {
	bytes, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key: %w", err)
	}
	vBytes, err := value.Serialize()
	if err != nil {
		return fmt.Errorf("serialize value: %w", err)
	}
	txn.MapKVClient[string(bytes)] = vBytes
	return nil
}

func (txn *MapKVTxn) Get(ctx context.Context, key Serializable, value Serializable) error {
	bytes, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key: %w", err)
	}

	vBytes, ok := txn.MapKVClient[string(bytes)]
	if !ok {
		return ErrRecordNotFound
	}

	err = value.UnSerialize(vBytes)
	if err != nil {
		return fmt.Errorf("unserialize value: %w", err)
	}

	return nil
}

func (txn *MapKVTxn) Delete(key Serializable) (err error) {
	bytes, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key: %w", err)
	}
	delete(txn.MapKVClient, string(bytes))
	return nil
}

func (client *MapKVClient) Begin() (KVTxn, error) {
	return &MapKVTxn{*client}, nil
}

func (client *MapKVClient) Close() error {
	return nil
}

func (txn *MapKVTxn) Commit(ctx context.Context) (err error) {
	return nil
}

func (txn *MapKVTxn) Rollback() (err error) {
	return nil
}
