package userfile

import (
	"github.com/Littlefisher619/cosdisk/repository/connections"
)

func NewMap() *UserfileStorageKV {
	return &UserfileStorageKV{
		userfileKV: connections.InitMapKV(),
	}
}
