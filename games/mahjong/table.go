package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"
	"globaltedinc/framework/network"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	MSG "globaltedinc/mahjongtester/proto"
)

type tableUserStatusT int32

const (
	playerStatusEmpty = tableUserStatusT(iota)
	playerStatusLogining
	playerStatusPrepairing
	playerStatusTakingIn
	playerStatusGetReady
	playerStatusBusy
	playerStatusWaiting
	playerStatusActive
	playerStatusChoosing
	playerStatusOffline
)

type tableUserT struct {
	uin    uint32
	name   string
	seat   int32
	status tableUserStatusT
	money  int32
	icon   string
}

type tableT struct {
	roomID   int32
	selfSeat int32

	players []tableUserT

	game *MahjongClientT
	user *playerDataT

	network         network.TCPClient
	fightServerAddr string

	stateManager base.StateManager
}

func (table *tableT) initialize() {
	table.players = make([]tableUserT, 2)

	table.stateManager.Initialize()

	stateSettings := []statePair{
		{tableStateIDConnectFightServer, &tableStateConnectFightServer{tableStateBase{table: table, stateManager: &table.stateManager, user: table.user}}},
		{tableStateIDLoginFight, &tableStateLoginFight{tableStateBase{table: table, stateManager: &table.stateManager, user: table.user}}},
		{tableStateIDTakeIn, &tableStateTakeIn{tableStateBase{table: table, stateManager: &table.stateManager, user: table.user}}},
		{tableStateIDPlayerReady, &tableStatePlayerReady{tableStateBase{table: table, stateManager: &table.stateManager, user: table.user}}},
	}

	// 注册状态
	for state := range stateSettings {
		if err := table.stateManager.RegisterState(stateSettings[state].id, stateSettings[state].state); err != nil {
			glog.Error("error:", err, "state id:", stateSettings[state].id)
			return
		}
	}

	// 设置初始状态
	table.stateManager.PushState(tableStateIDConnectFightServer, false)
}

func getPlayerSeatIndex(seat int32) int32 {
	return (seat - 1) / 2
}

func (table *tableT) addPlayer(uin uint32, name string, seat int32) {
	table.players[getPlayerSeatIndex(seat)] = tableUserT{uin: uin, name: name, seat: seat}
}

func (table *tableT) connect(addr string) error {
	glog.Info("connecting server:", addr)
	if err := table.network.Connect(addr, 10000,
		table.disconnectHandler,

		func(packet *network.Packet) {
			table.stateManager.SendMessage(packet.GetData())
		}); err != nil {
		return err
	}
	return nil
}

func (table *tableT) disconnect() {
	table.network.Disconnect()
}

func (table *tableT) disconnectHandler(addr string, err error) {
	glog.Infof("user %d, Server %s disconnected", table.user.uid, addr)
}

func (table *tableT) sendMessage(msg *MSG.LotillMsg) {
	if table.user.uin != 0 {
		msg.Uin = proto.Int32(int32(table.user.uin))
	}
	msg.Account = proto.String(fmt.Sprintf("%d", table.user.uid))
	msg.Name = proto.String(table.user.name)
	msg.TimeStamp = proto.Int64(table.user.timestamp)
	msg.Extra = proto.String(table.user.extra)
	platform := MSG.LoginPlatform_LOGIN_PLATFORM_IOS
	msg.Platform = &platform
	msg.RoomId = proto.Int32(table.roomID)

	data, err := proto.Marshal(msg)
	if err != nil {
		glog.Error("proto marchal error:", err, ". msg:", msg)
		return
	}
	//glog.Info("send data:", data)

	var p network.Packet
	p.Attach(data)
	table.network.SendPacket(&p)
}

func (table *tableT) update() {
	table.stateManager.Update()
}
