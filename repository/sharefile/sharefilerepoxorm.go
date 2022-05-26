package sharefile

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/pkg/dbdriver"
)

type SharefileStorageXORM struct {
	conn *dbdriver.DBEngine
}

func NewXORM(DBEngine *dbdriver.DBEngine) *SharefileStorageXORM {
	return &SharefileStorageXORM{
		conn: DBEngine,
	}
}

func (ss *SharefileStorageXORM) CreateShareFile(ctx context.Context, UserId string, Path string, ExpireDuration time.Duration) (string, error) {
	var ExpireTime time.Time
	if ExpireDuration == 0 {
		ExpireTime = time.Time{}
	} else {
		ExpireTime = time.Now().Add(ExpireDuration)
	}
	s := &model.ShareFile{
		UserId:     UserId,
		Path:       Path,
		ExpireTime: ExpireTime,
	}
	ss.conn.Db.Insert(s)
	return fmt.Sprint(s.ShareId), nil
}

func (ss *SharefileStorageXORM) GetVaildShareFile(ctx context.Context, ShareId string) (*model.ShareFile, error) {
	i, err := strconv.Atoi(ShareId)
	if err != nil {
		return nil, err
	}
	sf := &model.ShareFile{ShareId: int64(i)}
	has, err := ss.conn.Db.Get(sf)
	if err != nil || !has {
		return nil, model.ErrShareIdNotFound
	}
	if (sf.ExpireTime == time.Time{}) {
		return sf, nil
	}
	if time.Now().After(sf.ExpireTime) {
		return nil, model.ErrShareExpired
	}
	return sf, nil
}
