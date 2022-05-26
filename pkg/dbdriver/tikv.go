package dbdriver

import (
	"context"
	"fmt"

	tikverr "github.com/tikv/client-go/v2/error"
	"github.com/tikv/client-go/v2/txnkv"
	"github.com/tikv/client-go/v2/txnkv/transaction"
)

type TikvClient struct {
	*txnkv.Client
}

type TikvTXN struct {
	*transaction.KVTxn
}

func InitTikv(pdAddrs []string) *TikvClient {
	client, err := txnkv.NewClient(pdAddrs)
	if err != nil {
		fmt.Print("connect tikv client failed!")
		return nil
	}
	return &TikvClient{client}
}

func (txn *TikvTXN) Set(key Serializable, value Serializable) (err error) {
	keyBytes, err := key.Serialize()
	if err != nil {
		return err
	}
	valueBytes, err := value.Serialize()
	if err != nil {
		return err
	}
	err = txn.KVTxn.Set(keyBytes, valueBytes)
	if err != nil {
		return err
	}
	return nil
}

func (txn *TikvTXN) Get(ctx context.Context, key Serializable, value Serializable) (err error) {
	keyBytes, err := key.Serialize()
	if err != nil {
		return err
	}

	valueBytes, err := txn.KVTxn.Get(ctx, keyBytes)
	if tikverr.IsErrNotFound(err) {
		return ErrRecordNotFound
	}
	return value.UnSerialize(valueBytes)
}

func (txn *TikvTXN) Delete(key Serializable) (err error) {
	keyBytes, err := key.Serialize()
	if err != nil {
		return err
	}

	return txn.KVTxn.Delete(keyBytes)
}

func (rc *TikvClient) Begin() (KVTxn, error) {
	txn, err := rc.Client.Begin()
	if err != nil {
		return nil, err
	}
	return &TikvTXN{txn}, nil
}

var _ KVStorage = (*TikvClient)(nil)
