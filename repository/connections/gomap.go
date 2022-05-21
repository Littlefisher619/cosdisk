package connections

import (
	"github.com/Littlefisher619/cosdisk/model"
)

type MapkvConfig struct {
	m map[string]string
}

type MapkvTXN struct {
	m *MapkvConfig
}

func InitMapKV() *MapkvConfig {
	return &MapkvConfig{
		m: map[string]string{},
	}
}

func (rc *MapkvTXN) Set(key string, value string) (err error) {
	rc.m.m[key] = value
	return nil
}

func (rc *MapkvTXN) Get(key string) (value string, err error) {
	value, ok := rc.m.m[key]
	if !ok {
		return "", model.ErrFileNotFound
	}
	//fmt.Println("get ", string([]byte(key)), string(v))
	return value, nil
}

func (rc *MapkvTXN) Delete(key string) (err error) {
	delete(rc.m.m, key)
	return nil
}

func (rc *MapkvConfig) StartTranscation() (model.KeyValueTXN, error) {
	return &MapkvTXN{rc}, nil
}

func (rc *MapkvTXN) CommitTranscation() (err error) {
	return nil
}

func (rc *MapkvTXN) RollingBackTranscation() (err error) {
	return nil
}
