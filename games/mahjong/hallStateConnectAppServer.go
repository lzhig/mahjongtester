package mahjong

import (
	"globaltedinc/framework/base"

	"github.com/golang/glog"
)

type hallStateConnectAppServer struct {
	hallStateBase
}

func (state hallStateConnectAppServer) GetStateID() base.StateID {
	return hallStateIDConnectAppServer
}

func (state hallStateConnectAppServer) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", connecting app server...")
	// 连接服务器
	if err := state.hall.connect(state.hall.appServerAddr); err != nil {
		glog.Error("user:", state.user.uid, ", error: failed to connect app server.")
		return
	}
	glog.Info("user:", state.user.uid, ", connected to app server.")
	state.stateManager.ChangeState(hallStateIDLoginRole, false)
}
