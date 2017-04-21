package mahjong

import (
	"globaltedinc/framework/base"
	"time"

	MSG "globaltedinc/mahjongtester/proto"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type tableStatePlayerReady struct {
	tableStateBase
}

func (state tableStatePlayerReady) GetStateID() base.StateID {
	return tableStateIDPlayerReady
}

func (state tableStatePlayerReady) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", player ready...")

	state.request()
}

func (state tableStatePlayerReady) request() {
	msg := &MSG.LotillMsg{
		Msgid: proto.Int32(int32(MSG.MessageID_MSGID_PLAYER_READY_REQUEST)),
		PlayerReadyRequest: &MSG.PlayerReadyRequest{
			Ready: proto.Int32(4),
		},
	}

	state.table.sendMessage(msg)
}

func (state tableStatePlayerReady) OnMessage(data []byte) {
	msg := &MSG.LotillMsg{}
	if err := proto.Unmarshal(data, msg); err != nil {
		glog.Error("error:", err)
		return
	}

	glog.Info(msg)

	switch MSG.MessageID(*msg.Msgid) {
	case MSG.MessageID_MSGID_CALCULATE_NOTIFY:
		time.Sleep(time.Second * 5)
		state.request()
		glog.Info("Player Ready request.")
		//state.stateManager.ChangeState(hallStateIDReconnectInfo, false)
		//glog.Info("ChangeState(hallStateIDReconnectInfo, false)")

	case MSG.MessageID_MSGID_PLAYER_CHOICE_NOTIFY:
		if *msg.PlayerChoiceNotify.ChoiceType == MSG.ChoiceType_DA {
			time.Sleep(time.Second)
			choiceType := MSG.ChoiceType_ZHUA
			msg := &MSG.LotillMsg{
				Msgid: proto.Int32(int32(MSG.MessageID_MSGID_PLAYER_CHOICE_REQUEST)),
				PlayerChoiceRequest: &MSG.PlayerChoiceRequest{
					ChoiceType: &choiceType,
				},
			}
			state.table.sendMessage(msg)
		}
	case MSG.MessageID_MSGID_PLAYER_CHOICE_RESPONSE:
		if *msg.PlayerChoiceResponse.ChoiceType == MSG.ChoiceType_ZHUA {
			time.Sleep(time.Second * 2)
			choiceType := MSG.ChoiceType_DA
			msg := &MSG.LotillMsg{
				Msgid: proto.Int32(int32(MSG.MessageID_MSGID_PLAYER_CHOICE_REQUEST)),
				PlayerChoiceRequest: &MSG.PlayerChoiceRequest{
					ChoiceType: &choiceType,
					Mj:         []*MSG.MahjongDef{msg.PlayerChoiceResponse.Mj},
				},
			}
			state.table.sendMessage(msg)
		}

	case MSG.MessageID_MSGID_PLAYER_READY_RESPONSE:
		glog.Info("Player ready.")

	case MSG.MessageID_MSGID_SHUFFERING_NOTIFY:
		glog.Info("shuffering notify.")
		if *msg.ShufferingNotify.Banker == state.table.selfSeat {
			time.Sleep(time.Second * 2)
			choiceType := MSG.ChoiceType_DA
			msg := &MSG.LotillMsg{
				Msgid: proto.Int32(int32(MSG.MessageID_MSGID_PLAYER_CHOICE_REQUEST)),
				PlayerChoiceRequest: &MSG.PlayerChoiceRequest{
					ChoiceType: &choiceType,
					Mj:         []*MSG.MahjongDef{msg.ShufferingNotify.Mj[0]},
				},
			}
			state.table.sendMessage(msg)
		}

	case MSG.MessageID_MSGID_PLAYER_LOGIN_NOTIFY:
		glog.Info("player login notify.")

	case MSG.MessageID_MSGID_ERROR_CODE_NOTIFY:
		glog.Info("Error Code:", *msg.ErrorcodeNotify.ErrorCode, ", msg:", msg)

	default:
		glog.Error("Invalid message:", msg)
		return
	}
}
