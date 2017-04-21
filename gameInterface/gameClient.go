package gameInterface

import (
	"globaltedinc/mahjongtester/platform"
)

type AppID string

type GameClientInterface interface {
	GetAppID() AppID
	Start(uid uint32, name string, timestamp int64, extra string, token string, user *platform.UserT) error
}
