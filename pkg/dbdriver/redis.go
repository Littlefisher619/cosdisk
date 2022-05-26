package dbdriver

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
)

type Redis struct {
	*redis.Client
}

type RedisTxn struct {
	m *Redis
}

func InitRedis(addr string, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Print("connect redis failed!")
		return nil
	}
	return &Redis{
		client,
	}
}

func (txn *RedisTxn) Set(key Serializable, value Serializable) error {
	keyBytes, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key failed: %v", err)
	}
	valueBytes, err := value.Serialize()
	if err != nil {
		return fmt.Errorf("serialize value failed: %v", err)
	}
	resp := txn.m.Set(string(keyBytes), string(valueBytes), 0)
	return resp.Err()
}

func (rc *RedisTxn) Get(ctx context.Context, key Serializable, value Serializable) error {
	keyBytes, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key failed: %v", err)
	}
	valueBytes, err := rc.m.Get(string(keyBytes)).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrRecordNotFound
		}
		return fmt.Errorf("get value failed: %v", err)
	}
	err = value.UnSerialize([]byte(valueBytes))
	if err != nil {
		return fmt.Errorf("unserialize value failed: %v", err)
	}
	return nil
}

func (rc *RedisTxn) Delete(key Serializable) error {
	keyBytes, err := key.Serialize()
	if err != nil {
		return fmt.Errorf("serialize key failed: %v", err)
	}
	resp := rc.m.Del(string(keyBytes))
	return resp.Err()
}

func (rc *Redis) Begin() (KVTxn, error) {
	return &RedisTxn{rc}, nil
}

func (rc *RedisTxn) Commit(ctx context.Context) (err error) {
	return nil
}

func (rc *RedisTxn) Rollback() (err error) {
	return nil
}
