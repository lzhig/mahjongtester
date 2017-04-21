package mahjong

import (
	"fmt"
	"globaltedinc/framework/base"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type tableStateTakeIn struct {
	tableStateBase
}

func (state tableStateTakeIn) GetStateID() base.StateID {
	return tableStateIDTakeIn
}

func (state tableStateTakeIn) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", take in...")

	cash, _, _, err := state.table.game.userPlatform.GetBalance()
	if err != nil {
		glog.Error("error:", err)
		return
	}

	glog.Info("cash:", cash)

	response, err := state.table.game.userPlatform.AssetsIn(string(state.table.game.GetAppID()), "cash", cash, uint64(state.table.roomID), "room", state.user.token)
	if err != nil {
		glog.Error("error:", err)
		return
	}

	//glog.Info(response)

	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_TAKE_IN_REQUEST)),
		CashTakeInRequest: &MSG.CashTakeInRequest{
			Amount:   proto.String(fmt.Sprintf("%d", response.Amount)),
			AppId:    proto.String(response.AppID),
			Currency: proto.String(response.Currency),
			RoomId:   proto.String(fmt.Sprintf("%d", state.table.roomID)),
			Token:    proto.String(response.Token),
			Uid:      proto.String(fmt.Sprintf("%d", response.UID)),
		},
	}

	glog.Info("msg:", msg)

	state.table.sendMessage(msg)
}

func (state tableStateTakeIn) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	if *msg.Msgid != int32(MSG.MessageID_MSGID_TAKE_IN_RESPONSE) {
		glog.Error("Invalid message:", msg)
		return
	}

	glog.Info(msg)

	if *msg.CashTakeInResponse.Result == 0 {
		state.stateManager.ChangeState(tableStateIDPlayerReady, false)
	} else {
		glog.Warning("error: Failed to take in.")
		state.stateManager.ChangeState(tableStateIDLoginFight, false)
	}
}
