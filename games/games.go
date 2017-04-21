package games

import (
	"globaltedinc/mahjongtester/gameInterface"
	"globaltedinc/mahjongtester/games/mahjong"
)

func CreateGameClient(appid gameInterface.AppID) gameInterface.GameClientInterface {
	switch appid {
	case "241d7b34d61b3aea4eb0e3a91d753ecb":
		return &mahjong.MahjongClientT{Appid: appid}
	default:
		return nil
	}
}
