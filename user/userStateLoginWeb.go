package user

import (
	"globaltedinc/framework/base"

	"github.com/golang/glog"
)

type userStateLoginWeb struct {
	hallStateBase
}

func (state userStateLoginWeb) GetStateID() base.StateID {
	return stateIDLoginWeb
}

func (state userStateLoginWeb) OnEnter(prevStateID base.StateID) {
	glog.Info("user", state.user.uid, "login on web...")

	if err := state.user.LoginPlatform(state.user.username, state.user.password); err != nil {
		glog.Error("error:", err)
		return
	}

	state.stateManager.ChangeState(hallStateIDConnectIndexServer, false)
}
