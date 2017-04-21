package mahjong

import "globaltedinc/framework/base"

const (
	hallStateIDConnectIndexServer = base.StateID(base.StateInvalidID + iota + 1)
	hallStateIDListRole
	hallStateIDCreateRole
	hallStateIDConnectAppServer
	hallStateIDLoginRole
	hallStateIDReconnectInfo
	hallStateIDListRoom
	hallStateIDLoginRoom
	hallStateIDQuickLoginRoom
	hallStateIDCount
)

type hallStateBase struct {
	hall         *hallT
	user         *playerDataT
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
