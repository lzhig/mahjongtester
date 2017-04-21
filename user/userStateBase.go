package user

import "globaltedinc/framework/base"

const (
	stateIDLoginWeb = base.StateID(base.StateInvalidID + iota + 1)
	hallStateIDConnectIndexServer
	hallStateIDListRole
	hallStateIDCreateRole
	hallStateIDConnectAppServer
	hallStateIDLoginRole
	hallStateIDReconnectInfo
	hallStateIDListRoom
	hallStateIDLoginRoom
	hallStateIDQuickLoginRoom
	tableStateIDConnectFightServer
	tableStateIDLoginFight
	stateIDGetUserBalance
	stateIDAssetsIn
	tableStateIDTakeIn
	tableStateIDPlayerReady
)

type hallStateBase struct {
	user         *UserT
	stateManager *base.StateManager
}

func (state hallStateBase) OnEnter(prevStateID base.StateID) {

}

func (state hallStateBase) OnExit(nextStateID base.StateID) {

}

func (state hallStateBase) Update() {

}

func (state hallStateBase) OnMessage(data []byte) {

}

func (state hallStateBase) OnSuspend(nextStateID base.StateID) {

}

func (state hallStateBase) OnResume(prevStateID base.StateID) {

}
