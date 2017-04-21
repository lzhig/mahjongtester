package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateQuickLoginRoom struct {
	hallStateBase
}

func (state hallStateQuickLoginRoom) GetStateID() base.StateID {
	return hallStateIDQuickLoginRoom
}

func (state hallStateQuickLoginRoom) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", quick login room...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_QUICK_LOGIN_REQUEST)),
		QuickLoginRequest: &MSG.QuickLoginRequest{
			Grade: proto.Int32(1),
			//AnteType:     proto.Int32(state.hall.roomJoining.anteType),
			CurrencyType: proto.Int32(2),
		},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateQuickLoginRoom) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_QUICK_LOGIN_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	glog.Info(msg)

	if *msg.QuickLoginResponse.Result == 0 {
		state.hall.roomJoining.anteType = *msg.QuickLoginResponse.AnteType
		state.hall.roomJoining.id = *msg.QuickLoginResponse.RoomId
		state.hall.roomJoining.seat = *msg.QuickLoginResponse.SeatLocation

		glog.Info("quick login roomid:", state.hall.roomJoining.id, ", seat:", state.hall.roomJoining.seat)

		state.hall.game.createTable(fmt.Sprintf("%s:%s", *msg.QuickLoginResponse.Ip, *msg.QuickLoginResponse.Port), state.hall.roomJoining.id, state.hall.roomJoining.seat)
	} else {
		glog.Warning("Failed to quick login room.")
		state.stateManager.ChangeState(hallStateIDListRoom, false)
	}
}
