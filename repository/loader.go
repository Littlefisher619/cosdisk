package repository

import (
	"fmt"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/repository/account"
	"github.com/Littlefisher619/cosdisk/repository/connections"
	"github.com/Littlefisher619/cosdisk/repository/sharefile"
	"github.com/Littlefisher619/cosdisk/repository/userfile"
	"github.com/pelletier/go-toml"
)

type RepositoryConfig struct {
	AccountRepository   model.AccountRepository
	UserfileRepository  model.UserfileRepository
	ShareFileRepository model.ShareFileRepository
}

func (c *RepositoryConfig) CreateXormDB(driverName string, dataSourceName string) error {
	engine, err := connections.CreateXormEngine(driverName, dataSourceName, []interface{}{
		&model.User{},
		&model.ShareFile{},
	}, false)
	if err != nil {
		return err
	}
	c.ShareFileRepository = sharefile.NewXORM(engine)
	c.AccountRepository = account.NewXORM(engine)
	return nil
}

func LoadDataBaseConfig(config *toml.Tree) (RepositoryConfig, error) {
	db := config.GetDefault("database.sqlDB", "memmap")
	kv := config.GetDefault("database.keyValueDB", "memmap")
	c := RepositoryConfig{}

	switch kv {
	case "memmap":
		c.UserfileRepository = userfile.NewMap()
	case "redis":
		pass := config.GetDefault("redis.password", "")
		addr := config.GetDefault("redis.addr", "localhost:6379")
		redis := connections.InitRedis(addr.(string), pass.(string))
		if redis == nil {
			return c, fmt.Errorf("redis init failed")
		}
		c.UserfileRepository = userfile.NewUserfileKV(redis)
	case "tikv":
		pdaddr := config.GetDefault("tikv.pdAddr", "127.0.0.1:2379")
		tikv := connections.InitTikv([]string{pdaddr.(string)})
		c.UserfileRepository = userfile.NewUserfileKV(tikv)
	default:
		c.UserfileRepository = userfile.NewMap()
	}

	switch db {
	case "memmap":
		c.ShareFileRepository = sharefile.New()
		c.AccountRepository = account.NewMap()
	case "tidb":
		dataSourceName := config.GetDefault("tidb.dataSourceName", "root@tcp(127.0.0.1:4000)/test?charset=utf8")
		err := c.CreateXormDB("mysql", dataSourceName.(string))
		if err != nil {
			return c, err
		}
	default:
		c.ShareFileRepository = sharefile.New()
		c.AccountRepository = account.NewMap()
	}
	return c, nil
}
