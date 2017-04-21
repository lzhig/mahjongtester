package user

import (
	"globaltedinc/framework/base"

	"github.com/golang/glog"
)

type userStateGetUserBalance struct {
	hallStateBase
}

func (state userStateGetUserBalance) GetStateID() base.StateID {
	return stateIDGetUserBalance
}

func (state userStateGetUserBalance) OnEnter(prevStateID base.StateID) {
	glog.Info("user", state.user.uid, "get user balance...")

	if err := state.user.getBalance(); err != nil {
		glog.Error("error:", err)
		return
	}

	state.stateManager.ChangeState(stateIDAssetsIn, false)
}
