package connections

import (
	"fmt"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/go-redis/redis"
)

type RedisConfig struct {
	client *redis.Client
}

type RediskvTXN struct {
	m *RedisConfig
}

func InitRedis(addr string, password string) *RedisConfig {
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
	return &RedisConfig{
		client: client,
	}

}

func (rc *RediskvTXN) Set(key string, value string) (err error) {
	err = rc.m.client.Set(key, value, 0).Err()
	return
}

func (rc *RediskvTXN) Get(key string) (value string, err error) {
	value, err = rc.m.client.Get(key).Result()
	return
}

func (rc *RediskvTXN) Delete(key string) (err error) {
	err = rc.m.client.Del(key).Err()
	return
}

func (rc *RedisConfig) StartTranscation() (model.KeyValueTXN, error) {
	return &RediskvTXN{rc}, nil
}

func (rc *RediskvTXN) CommitTranscation() (err error) {
	return nil
}

func (rc *RediskvTXN) RollingBackTranscation() (err error) {
	return nil
}
