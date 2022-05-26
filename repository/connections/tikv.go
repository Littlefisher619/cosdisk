package connections

import (
	"context"
	"fmt"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/tikv/client-go/v2/txnkv"
	"github.com/tikv/client-go/v2/txnkv/transaction"
)

type TikvConfig struct {
	client *txnkv.Client
}

type TikvTXN struct {
	txn *transaction.KVTxn
}

func InitTikv(pdAddrs []string) *TikvConfig {
	client, err := txnkv.NewClient(pdAddrs)
	if err != nil {
		fmt.Print("connect tikv client failed!")
		return nil
	}
	return &TikvConfig{
		client: client,
	}
}

func (rc *TikvConfig) CloseTikv() error {
	return rc.client.Close()
}

func (rc *TikvTXN) Set(key string, value string) (err error) {
	err = rc.txn.Set([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	//fmt.Println("set ", string([]byte(key)), value)
	return nil
}

func (rc *TikvTXN) Get(key string) (value string, err error) {
	v, err := rc.txn.Get(context.Background(), []byte(key))
	if err != nil {
		return "", err
	}
	//fmt.Println("get ", string([]byte(key)), string(v))
	return string(v), nil
}

func (rc *TikvTXN) Delete(key string) (err error) {
	return rc.txn.Delete([]byte(key))
}

func (rc *TikvConfig) StartTranscation() (model.KeyValueTXN, error) {
	txn, err := rc.client.Begin()
	if err != nil {
		return nil, err
	}
	return &TikvTXN{txn}, nil
}

func (rc *TikvTXN) CommitTranscation() (err error) {
	return rc.txn.Commit(context.Background())
}

func (rc *TikvTXN) RollingBackTranscation() (err error) {
	return rc.txn.Rollback()
}
