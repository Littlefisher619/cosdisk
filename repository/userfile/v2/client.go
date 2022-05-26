package userfile

import (
	"context"
	"fmt"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/pkg/dbdriver"
)

type Client struct {
	dbdriver.KVStorage
}

type Txn struct {
	dbdriver.KVTxn
}

func NewClient(kv dbdriver.KVStorage) *Client {
	return &Client{
		KVStorage: kv,
	}
}

func (c *Client) Begin() (txn model.UserfileTXN, err error) {
	tn, err := c.KVStorage.Begin()
	if err != nil {
		return nil, err
	}
	return &Txn{tn}, nil
}

func (c *Client) RunInTxn(ctx context.Context, txnfunc model.TxnFunc) error {
	txn, err := c.Begin()
	if err != nil {
		return err
	}

	err = txnfunc(txn)
	if err != nil {
		if e := txn.Rollback(); e != nil {
			return fmt.Errorf("Rollback failed %s: %w", err, e)
		}
		return err
	}

	err = txn.Commit(ctx)
	if err != nil {
		if e := txn.Rollback(); e != nil {
			return e
		}
		return err
	}
	return nil
}

func (txn *Txn) Commit(ctx context.Context) (err error) {
	return txn.KVTxn.Commit(ctx)
}

func (txn *Txn) Rollback() (err error) {
	return txn.KVTxn.Rollback()
}
