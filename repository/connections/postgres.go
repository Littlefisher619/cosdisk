package connections

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type DBPgConn struct {
	Db *pg.DB
}

func NewPostgresDB(models []interface{}, isTemp bool) (*DBPgConn, error) {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
	})

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return nil, err
		}
	}
	return &DBPgConn{Db: db}, nil
}
