package dbdriver

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
)

type DBEngine struct {
	Db *xorm.Engine
}

func CreateXormEngine(driverName string, dataSourceName string, models []interface{}, isTemp bool) (*DBEngine, error) {
	// "root@/tcp(127.0.0.1:4000)"
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql failed %e", err)
	}
	if err := engine.Ping(); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("connect to mysql failed %e", err)
	}
	engine.ShowSQL(true)
	for _, model := range models {
		err := engine.Sync(model)
		if err != nil {
			return nil, err
		}
	}

	return &DBEngine{engine}, nil
}
