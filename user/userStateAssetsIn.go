package user

import (
	"globaltedinc/framework/base"

	"github.com/golang/glog"
)

type userStateAssetsIn struct {
	hallStateBase
}

func (state userStateAssetsIn) GetStateID() base.StateID {
	return stateIDAssetsIn
}

func (state userStateAssetsIn) OnEnter(prevStateID base.StateID) {
	glog.Info("user", state.user.uid, "assetsin...")
	/*
		assetsInResponse, err := state.user.assetsIn("")
		if err != nil {
			glog.Error("error:", err)
			return
		}

		if tmp, err := strconv.ParseInt(assetsInResponse.RoomID, 10, 32); err != nil {
			glog.Error("error:", err)
			return
		} else {
			state.user.roomID = int32(tmp)
		}
	*/
	state.stateManager.ChangeState(tableStateIDTakeIn, false)
}
