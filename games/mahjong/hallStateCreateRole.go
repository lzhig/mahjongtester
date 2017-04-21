package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateCreateRole struct {
	hallStateBase
}

func (state hallStateCreateRole) GetStateID() base.StateID {
	return hallStateIDCreateRole
}

func (state hallStateCreateRole) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", create role list...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_CREATEROLE_REQUEST)),
		CreateRoleRequest: &MSG.CreateRoleRequest{
			Gender: proto.Int32(1),
			Avatar: proto.Int32(2),
			Name:   proto.String(state.user.name),
		},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateCreateRole) Update() {

}

func (state hallStateCreateRole) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_CREATEROLE_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	glog.Info(msg)
	if *msg.CreateRoleResponse.Result == 0 {
		ip := []byte(*msg.CreateRoleResponse.Ip)
		len := 0
		for ndx := range ip {
			if ip[ndx] == 0 {
				break
			}
			len++
		}

		state.user.uin = *msg.CreateRoleResponse.Uin
		state.hall.appServerAddr = fmt.Sprintf("%s:%d", string(ip[:len]), *msg.CreateRoleResponse.Port)

		state.hall.disconnect()
		state.stateManager.ChangeState(hallStateIDConnectAppServer, false)
	}
}
