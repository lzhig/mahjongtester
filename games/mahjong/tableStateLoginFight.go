package mahjong

import (
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type tableStateLoginFight struct {
	tableStateBase
}

func (state tableStateLoginFight) GetStateID() base.StateID {
	return tableStateIDLoginFight
}

func (state tableStateLoginFight) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", login fight...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_LOGIN_FIGHT_REQUEST)),
		LoginFightRequest: &MSG.LoginFightRequest{
			RoomId:       proto.Int32(state.table.roomID),
			SeatLocation: proto.Int32(state.table.selfSeat),
		},
	}

	state.table.sendMessage(msg)
}

func (state tableStateLoginFight) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_LOGIN_FIGHT_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	glog.Info(msg)

	if *msg.LoginFightResponse.Result > 0 {
		state.stateManager.ChangeState(tableStateIDTakeIn, false)
	} else {
		glog.Warning("error: Failed to login fight.")
		//state.stateManager.ChangeState(hallStateIDListRoom, false)
	}
}
