package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"
	"globaltedinc/framework/network"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"

	MSG "globaltedinc/mahjongtester/proto"
)

const (
	CashRoom  = iota + 1 // 现金场
	CoinRoom             // 虚拟币场
	MatchRoom            // 比赛场
)

type hallT struct {
	user         *playerDataT
	game         *MahjongClientT
	stateManager base.StateManager

	roomJoining roomT

	network         network.TCPClient
	indexServerAddr string
	appServerAddr   string
}

type statePair struct {
	id    base.StateID
	state base.StateInterface
}

func (hall *hallT) initialize() {
	hall.stateManager.Initialize()

	stateSettings := []statePair{
		{hallStateIDConnectIndexServer, &hallStateConnectIndexServer{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDCreateRole, &hallStateCreateRole{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDListRole, &hallStateListRole{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDConnectAppServer, &hallStateConnectAppServer{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDLoginRole, &hallStateLoginRole{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDReconnectInfo, &hallStateReconnectInfo{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDListRoom, &hallStateListRoom{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDLoginRoom, &hallStateLoginRoom{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
		{hallStateIDQuickLoginRoom, &hallStateQuickLoginRoom{hallStateBase{hall: hall, stateManager: &hall.stateManager, user: hall.user}}},
	}

	// 注册状态
	for state := range stateSettings {
		if err := hall.stateManager.RegisterState(stateSettings[state].id, stateSettings[state].state); err != nil {
			glog.Error("error:", err, "state id:", stateSettings[state].id)
			return
		}
	}

	// 设置初始状态
	hall.stateManager.PushState(hallStateIDConnectIndexServer, false)
}

func (hall *hallT) update() {
	hall.stateManager.Update()
}

func (hall *hallT) connect(addr string) error {
	glog.Info("connecting server:", addr)
	if err := hall.network.Connect(addr, 10000,
		hall.disconnectHandler,

		func(packet *network.Packet) {
			hall.stateManager.SendMessage(packet.GetData())
		}); err != nil {
		return err
	}
	return nil
}

func (hall *hallT) disconnect() {
	hall.network.Disconnect()
}

func (hall *hallT) disconnectHandler(addr string, err error) {
	glog.Infof("user %d, Server %s disconnected", hall.user.uid, addr)
}

func (hall *hallT) sendMessage(msg *MSG.LotillMsg) {
	if hall.user.uin != 0 {
		msg.Uin = proto.Int32(int32(hall.user.uin))
	}
	msg.Account = proto.String(fmt.Sprintf("%d", hall.user.uid))
	msg.Name = proto.String(hall.user.name)
	msg.TimeStamp = proto.Int64(hall.user.timestamp)
	msg.Extra = proto.String(hall.user.extra)
	platform := MSG.LoginPlatform_LOGIN_PLATFORM_IOS
	msg.Platform = &platform
	//msg.RoomId = proto.Int32(hall.user.roomID)

	data, err := proto.Marshal(msg)
	if err != nil {
		glog.Error("proto marchal error:", err, ". msg:", msg)
		return
	}
	//glog.Info("send data:", data)

	var p network.Packet
	p.Attach(data)
	hall.network.SendPacket(&p)
}
