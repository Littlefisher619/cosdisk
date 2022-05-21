package sharefile

import (
	"context"
	"fmt"
	"time"

	"github.com/Littlefisher619/cosdisk/model"
)

type SharefileStorageMap struct {
	ssmap map[string]*model.ShareFile
}

func New() *SharefileStorageMap {
	return &SharefileStorageMap{
		ssmap: make(map[string]*model.ShareFile),
	}
}

// need database or lock
// for test only
var id int64 = 0

func (ss *SharefileStorageMap) CreateShareFile(ctx context.Context, UserId string, Path string, ExpireDuration time.Duration) (string, error) {
	var ExpireTime time.Time
	if ExpireDuration == 0 {
		ExpireTime = time.Time{}
	} else {
		ExpireTime = time.Now().Add(ExpireDuration)
	}
	ss.ssmap[fmt.Sprint(id)] = &model.ShareFile{
		ShareId:    id,
		UserId:     UserId,
		Path:       Path,
		ExpireTime: ExpireTime,
	}
	id = id + 1
	return fmt.Sprint(id - 1), nil
}

func (ss *SharefileStorageMap) GetVaildShareFile(ctx context.Context, ShareId string) (*model.ShareFile, error) {
	sf, ok := ss.ssmap[ShareId]
	if !ok {
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
