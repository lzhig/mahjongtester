package user

import (
	"fmt"
	"globaltedinc/framework/base"
	"globaltedinc/framework/network"
	"globaltedinc/mahjongtester/gameInterface"
	"globaltedinc/mahjongtester/platform"
	MSG "globaltedinc/mahjongtester/proto"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

const (
	CashRoom  = iota + 1 // 现金场
	CoinRoom             // 虚拟币场
	MatchRoom            // 比赛场
)

type UserT struct {
	userPlatformT

	gamesPlaying map[string]gameInterface.GameClientInterface

	roomAnte     int32  // 房间底注
	roomID       int32  // 房间ID
	roomSeat     int32  // 房间位置
	roomAnteType int32  // 房间类型
	roomCurrency string // 货币类型

	fightServer string

	uin           uint32
	appServer     string
	appServerPort int32

	IndexServer string
	IndexNet    network.TCPClient
	AppNet      network.TCPClient
	fightNet    network.TCPClient

	assetsInResponse assetsInResponseT

	//data chan *network.Packet

	stateManager base.StateManager
}

func (user *UserT) Initialize(platform *platform.PlatformT) error {
	user.uin = 0
	//user.data = make(chan *network.Packet, 1)
	user.roomCurrency = "cash"

	user.gamesPlaying = make(map[string]gameInterface.GameClientInterface)

	if err := user.userPlatformT.initialize(platform); err != nil {
		return err
	}

	return nil
}

func (user *UserT) setAppServerInfo(uin uint32, server string, port int32) {
	user.uin = uin
	user.appServer = server
	user.appServerPort = port
}

type statePair struct {
	id    base.StateID
	state base.StateInterface
}

func (user *UserT) StartGame(wg *sync.WaitGroup) {

	user.stateManager.Initialize()

	stateSettings := []statePair{
		{stateIDLoginWeb, &userStateLoginWeb{hallStateBase{user: user, stateManager: &user.stateManager}}},

		{stateIDGetUserBalance, &userStateGetUserBalance{hallStateBase{user: user, stateManager: &user.stateManager}}},
		{stateIDAssetsIn, &userStateAssetsIn{hallStateBase{user: user, stateManager: &user.stateManager}}},
		
		
	}

	// 注册状态
	for state := range stateSettings {
		if err := user.stateManager.RegisterState(stateSettings[state].id, stateSettings[state].state); err != nil {
			glog.Error("error:", err, "state id:", stateSettings[state].id)
			return
		}
	}

	// 设置初始状态
	user.stateManager.PushState(stateIDLoginWeb, false)

	go user.update()
}

func (user *UserT) EndGame(wg *sync.WaitGroup) error {
	return nil
}

func (user *UserT) connectIndexServer() error {
	glog.Info("index server:", user.IndexServer)
	if err := user.IndexNet.Connect(user.IndexServer, 10000,
		func(addr string, err error) {
			glog.Error("user", user.uid, "Index Server disconnected.")
		},

		func(packet *network.Packet) {
			//glog.Info("received data...")
			//data := packet.GetData()
			user.stateManager.SendMessage(packet.GetData())
			//lenBuf := bytes.NewBuffer([]byte{})
			//binary.Write(lenBuf, binary.BigEndian, len(data))
			//glog.Info("len(data)=", len(data))
			//glog.Info("data:", data)
			//user.data <- packet
		}); err != nil {
		return err
	}
	return nil
}

func (user *UserT) connectAppServer() error {
	if err := user.AppNet.Connect(fmt.Sprintf("%s:%d", user.appServer, user.appServerPort), 10000,
		func(addr string, err error) {
			glog.Error("user", user.uid, "disconnected. reconnecting app server...")
			user.connectAppServer()
		},

		func(packet *network.Packet) {
			//glog.Info("received data...")
			user.stateManager.SendMessage(packet.GetData())
			//glog.Info("len(data)=", len(packet.GetData()))
		}); err != nil {
		return err
	}

	return nil
}

func (user *UserT) connectFightServer() error {
	if err := user.fightNet.Connect(user.fightServer, 10000,
		func(addr string, err error) {
			glog.Error("user", user.uid, "disconnected. reconnecting fight server...")
			user.connectFightServer()
		},

		func(packet *network.Packet) {
			//glog.Info("received data...")
			user.stateManager.SendMessage(packet.GetData())
			//glog.Info("len(data)=", len(packet.GetData()))

		}); err != nil {
		return err
	}

	return nil
}

func (user *UserT) update() {
	t1 := time.Now()
	for {
		user.stateManager.Update()
		t2 := time.Now()
		delta := t2.Sub(t1)
		if delta < time.Millisecond {
			time.Sleep(time.Millisecond - delta)
		}
		t1 = t2
	}
}

func (user *UserT) sendMessage(net network.TCPClient, msg *MSG.LotillMsg) {
	if user.uin != 0 {
		msg.Uin = proto.Int32(int32(user.uin))
	} /*
		msg.Account = proto.String(fmt.Sprintf("%d", user.uid))
		msg.Name = proto.String(user.name)
		msg.TimeStamp = proto.Int64(user.timestamp)
		msg.Extra = proto.String(user.extra)
		platform := MSG.LoginPlatform_LOGIN_PLATFORM_IOS
		msg.Platform = &platform
		msg.RoomId = proto.Int32(user.roomID)
	*/
	data, err := proto.Marshal(msg)
	if err != nil {
		glog.Error("proto marchal error:", err)
		return
	}
	//glog.Info("send data:", data)

	var p network.Packet
	p.Attach(data)
	net.SendPacket(&p)
}
