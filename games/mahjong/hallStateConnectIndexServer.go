package mahjong

import (
	"globaltedinc/framework/base"

	"github.com/golang/glog"
)

type hallStateConnectIndexServer struct {
	hallStateBase
}

func (state hallStateConnectIndexServer) GetStateID() base.StateID {
	return hallStateIDConnectIndexServer
}

func (state hallStateConnectIndexServer) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", connecting index server...")
	// 连接服务器
	if err := state.hall.connect(state.hall.indexServerAddr); err != nil {
		glog.Error("user:", state.user.uid, ", error: failed to connect index server.")
		return
	}
	glog.Info("user:", state.user.uid, ", connected to index server.")

	state.stateManager.ChangeState(hallStateIDListRole, false)
}
