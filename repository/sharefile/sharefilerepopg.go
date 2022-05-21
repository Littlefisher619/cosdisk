package sharefile

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Littlefisher619/cosdisk/model"
	"github.com/Littlefisher619/cosdisk/repository/connections"
)

type SharefileStoragePG struct {
	conn *connections.DBPgConn
}

func NewPG(DBPgConn *connections.DBPgConn) *SharefileStoragePG {
	return &SharefileStoragePG{
		conn: DBPgConn,
	}
}

func (ss *SharefileStoragePG) CreateShareFile(ctx context.Context, UserId string, Path string, ExpireDuration time.Duration) (string, error) {
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
	ss.conn.Db.Model(&s).Insert()
	return fmt.Sprint(s.ShareId), nil
}

func (ss *SharefileStoragePG) GetVaildShareFile(ctx context.Context, ShareId string) (*model.ShareFile, error) {
	i, err := strconv.Atoi(ShareId)
	if err != nil {
		return nil, err
	}
	sf := &model.ShareFile{ShareId: int64(i)}
	err = ss.conn.Db.Model(sf).WherePK().Select()
	if err != nil {
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
