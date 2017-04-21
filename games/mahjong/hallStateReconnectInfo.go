package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateReconnectInfo struct {
	hallStateBase
}

func (state hallStateReconnectInfo) GetStateID() base.StateID {
	return hallStateIDReconnectInfo
}

func (state hallStateReconnectInfo) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", ReconnectInfo...")

	msg := &MSG.LotillMsg{
		Msgid:                proto.Int32(int32(MSG.MessageID_MSGID_RECONNECT_INFO_REQUEST)),
		ReconnectInfoRequest: &MSG.ReconnectInfoRequest{},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateReconnectInfo) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_RECONNECT_INFO_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	glog.Info(msg)

	if *msg.ReconnectInfoResponse.Result > 0 {
		state.hall.roomJoining.anteType = *msg.ReconnectInfoResponse.AnteType
		state.hall.roomJoining.seat = *msg.ReconnectInfoResponse.SeatLocation
		state.hall.roomJoining.id = *msg.ReconnectInfoResponse.RoomId

		state.hall.game.createTable(fmt.Sprintf("%s:%s", *msg.ReconnectInfoResponse.Ip, *msg.ReconnectInfoResponse.Port), state.hall.roomJoining.id, state.hall.roomJoining.seat)
	} else {
		glog.Warning("no Reconnect Info")
		state.stateManager.ChangeState(hallStateIDQuickLoginRoom, false)
	}
}
