package mahjong

import (
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type hallStateListRoom struct {
	hallStateBase
}

func (state hallStateListRoom) GetStateID() base.StateID {
	return hallStateIDListRoom
}

func (state hallStateListRoom) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", get room list...")

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_LIST_ROOM_REQUEST)),
		ListRoomRequest: &MSG.ListRoomRequest{
			RoomType:     proto.Int32(0),
			CurrencyType: proto.Int32(CashRoom),
		},
	}

	state.hall.sendMessage(msg)
}

func (state hallStateListRoom) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_LIST_ROOM_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	//glog.Info(msg)

	// 查找一个空桌
	for ndx := range msg.ListRoomResponse.RoomState {
		roomState := msg.ListRoomResponse.RoomState[ndx]
		playerNumber := len(roomState.PlayerState)
		if playerNumber == 0 {
			state.hall.roomJoining.ante = *roomState.Ante
			state.hall.roomJoining.id = *roomState.RoomId
			state.hall.roomJoining.seat = 1
			state.stateManager.ChangeState(hallStateIDLoginRoom, false)
			glog.Info(roomState)
			return
		} else if playerNumber == 1 {
			state.hall.roomJoining.ante = *roomState.Ante
			state.hall.roomJoining.id = *roomState.RoomId
			state.hall.roomJoining.seat = 4 - (*roomState.PlayerState[0].Location)
			state.stateManager.ChangeState(hallStateIDLoginRoom, false)
			glog.Info(roomState)
			return
		}
	}
	state.stateManager.ChangeState(hallStateIDListRoom, false)
}
