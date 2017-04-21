package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"globaltedinc/framework/network"
	"globaltedinc/mahjongtester/gameInterface"
	"globaltedinc/mahjongtester/games"
	"globaltedinc/mahjongtester/platform"
	"globaltedinc/mahjongtester/user"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/golang/glog"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	defer glog.Flush()

	startupConfig, err := loadStartupConfig()
	if err != nil {
		glog.Error(err)
		return
	}

	gs := platform.CreatePlatform(startupConfig.Platform.Name, startupConfig.Platform.URL)
	for ndx := range startupConfig.Games {
		gs.AddGame(startupConfig.Games[ndx].Appid, startupConfig.Games[ndx].Name)
	}

	var mahjongPacketHeader user.MahjongPacketHeader
	network.SetPacketHeader(&mahjongPacketHeader)

	// 开始压测
	//var wg sync.WaitGroup

	for i := uint64(startupConfig.Users.Start); i <= startupConfig.Users.End; i++ {
		// 创建平台用户
		user, err := gs.CreateUser(fmt.Sprintf("%s%d", startupConfig.Users.UsernamePrefix, i), startupConfig.Users.Password)
		if err != nil {
			glog.Error(err)
			return
		}

		// 登录平台
		if err := user.LoginPlatform(); err != nil {
			glog.Error(err)
			return
		}

		// 打开游戏
		for ndx := range startupConfig.GamesOpen {
			client := games.CreateGameClient(startupConfig.GamesOpen[ndx])
			uid, name, timestamp, extra, token, err := user.GetGameInfo(string(client.GetAppID()))
			if err != nil {
				glog.Error(err)
				return
			}
			client.Start(uid, name, timestamp, extra, token, user)
		}
	}

	fmt.Println("Press 'q' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('q')
}

type jsonPlatformConfigT struct {
	Name string `jsin:"name"`
	URL  string `json:"url"`
}

type jsonUsersConfigT struct {
	UsernamePrefix string `json:"username_prefix"`
	Start          uint64 `json:"start"`
	End            uint64 `json:"end"`
	Password       string `json:"password"`
}

type jsonGamesConfigT struct {
	Appid string `jsin:"appid"`
	Name  string `jsin:"name"`
}

type jsonStartupConfigT struct {
	Platform  jsonPlatformConfigT   `json:"platform"`
	Users     jsonUsersConfigT      `json:"users"`
	Games     []jsonGamesConfigT    `json:"games"`
	GamesOpen []gameInterface.AppID `json:"games-open"`
}

func loadStartupConfig() (*jsonStartupConfigT, error) {
	data, err := ioutil.ReadFile("startup.json")
	if err != nil {
		return nil, err
	}

	cfg := &jsonStartupConfigT{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
