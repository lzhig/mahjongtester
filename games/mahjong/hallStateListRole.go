package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateListRole struct {
	hallStateBase
}

func (state hallStateListRole) GetStateID() base.StateID {
	return hallStateIDListRole
}

func (state hallStateListRole) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", get role list...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_LISTROLE_REQUEST)),
		ListRoleRequest: &MSG.ListRoleRequest{
			Token: []byte(state.user.token),
			Name:  proto.String(state.user.name),
		},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateListRole) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_LISTROLE_RESPONSE) {
		glog.Error("Invalid message:", msg)
	}

	if *msg.ListRoleResponse.Result == 127 {
		state.stateManager.ChangeState(hallStateIDCreateRole, false)
	} else if *msg.ListRoleResponse.Result == 0 {
		ip := []byte(*msg.ListRoleResponse.Appip)
		len := 0
		for ndx := range ip {
			if ip[ndx] == 0 {
				break
			}
			len++
		}

		state.user.uin = *msg.ListRoleResponse.Uin
		state.hall.appServerAddr = fmt.Sprintf("%s:%d", string(ip[:len]), *msg.ListRoleResponse.Port)

		state.hall.disconnect()
		state.stateManager.ChangeState(hallStateIDConnectAppServer, false)
	} else {
		glog.Warning("ListRoleResponse 返回错误")
	}
}
