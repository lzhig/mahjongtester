package mahjong

import (
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateLoginRole struct {
	hallStateBase
}

func (state hallStateLoginRole) GetStateID() base.StateID {
	return hallStateIDLoginRole
}

func (state hallStateLoginRole) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", login role...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_LOGIN_REQUEST)),
		LoginRequest: &MSG.LoginRequest{
			Channel: proto.Int32(1),
		},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateLoginRole) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_LOGIN_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	if *msg.LoginResponse.Result == 0 {
		state.user.uin = *msg.LoginResponse.Roleinfo.Uin

		glog.Info("user uin:", state.user.uin)

		state.stateManager.ChangeState(hallStateIDReconnectInfo, false)
	} else {
		glog.Warning("error: Failed to login")
	}
}
