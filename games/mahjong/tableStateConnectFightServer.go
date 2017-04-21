package mahjong

import (
	"globaltedinc/framework/base"

	"github.com/golang/glog"
)

type tableStateConnectFightServer struct {
	tableStateBase
}

func (state tableStateConnectFightServer) GetStateID() base.StateID {
	return tableStateIDConnectFightServer
}

func (state tableStateConnectFightServer) OnEnter(prevStateID base.StateID) {
	glog.Info("user:", state.user.uid, ", connecting fight server...")
	// 连接服务器
	if err := state.table.connect(state.table.fightServerAddr); err != nil {
		glog.Error("user:", state.user.uid, ", error: failed to connect fight server.")
		return
	}
	glog.Info("user:", state.user.uid, ", connected to fight server.")
	state.stateManager.ChangeState(tableStateIDLoginFight, false)
}
