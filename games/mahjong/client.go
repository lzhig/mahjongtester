package mahjong

import (
	"encoding/json"
	"globaltedinc/mahjongtester/gameInterface"
	"globaltedinc/mahjongtester/platform"
	"io/ioutil"
	"time"
)

type MahjongClientT struct {
	Appid      gameInterface.AppID
	ServerAddr string

	user playerDataT

	hall   *hallT
	tables map[int32]*tableT

	userPlatform *platform.UserT
}

func (game *MahjongClientT) GetAppID() gameInterface.AppID {
	return game.Appid
}

type jsonConfigT struct {
	Server string `jsin:"server"`
}

func (game *MahjongClientT) loadConfig() (*jsonConfigT, error) {
	data, err := ioutil.ReadFile("mahjong.json")
	if err != nil {
		return nil, err
	}
	var cfg jsonConfigT
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (game *MahjongClientT) Start(uid uint32, name string, timestamp int64, extra string, token string, user *platform.UserT) error {
	game.user.uid = uid
	game.user.name = name
	game.user.timestamp = timestamp
	game.user.extra = extra
	game.user.token = token
	game.userPlatform = user

	cfg, err := game.loadConfig()
	if err != nil {
		return err
	}

	game.hall = &hallT{user: &game.user, indexServerAddr: cfg.Server, game: game}
	game.hall.initialize()

	game.tables = make(map[int32]*tableT)

	go game.update()

	return nil
}

func (game *MahjongClientT) createTable(addr string, roomID int32, roomSeat int32) {
	table := &tableT{fightServerAddr: addr, roomID: roomID, game: game, user: &game.user, selfSeat: roomSeat}
	table.initialize()
	table.addPlayer(game.user.uin, game.user.name, roomSeat)
	game.tables[roomID] = table
}

func (game *MahjongClientT) update() {
	t1 := time.Now()
	for {

		game.hall.update()

		for _, v := range game.tables {
			v.update()
		}

		t2 := time.Now()
		delta := t2.Sub(t1)
		if delta < time.Millisecond {
			time.Sleep(time.Millisecond - delta)
		}
		t1 = t2
	}
}
