package userfile

import (
	"github.com/Littlefisher619/cosdisk/pkg/dbdriver"
)

func NewMap() *UserfileStorageKV {
	return &UserfileStorageKV{
		userfileKV: dbdriver.InitMapKV(),
	}
}
