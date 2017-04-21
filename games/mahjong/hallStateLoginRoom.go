package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateLoginRoom struct {
	hallStateBase
}

func (state hallStateLoginRoom) GetStateID() base.StateID {
	return hallStateIDLoginRoom
}

func (state hallStateLoginRoom) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", login room...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_LOGIN_ROOM_REQUEST)),
		LoginRoomRequest: &MSG.LoginRoomRequest{
			RoomId:       proto.Int32(state.hall.roomJoining.id),
			SeatLocation: proto.Int32(state.hall.roomJoining.seat),
			AnteType:     proto.Int32(state.hall.roomJoining.anteType),
		},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateLoginRoom) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_LOGIN_ROOM_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	glog.Info(msg)

	if *msg.LoginRoomResponse.Result == 0 {
		state.hall.roomJoining.anteType = *msg.LoginRoomResponse.AnteType

		state.hall.game.createTable(fmt.Sprintf("%s:%s", *msg.LoginRoomResponse.Ip, *msg.LoginRoomResponse.Port), state.hall.roomJoining.id, state.hall.roomJoining.seat)
	} else {
		glog.Warning("Failed to login room.")
		state.stateManager.ChangeState(hallStateIDListRoom, false)
	}
}
