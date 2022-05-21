package model

import (
	"context"
	"time"
)

type ShareFile struct {
	ShareId    int64
	UserId     string
	Path       string
	ExpireTime time.Time
}

type ShareFileRepository interface {
	CreateShareFile(ctx context.Context, UserId string, Path string, ExpireDuration time.Duration) (string, error)
	GetVaildShareFile(ctx context.Context, ShareId string) (*ShareFile, error)
}
