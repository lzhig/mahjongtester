package mahjong

import "globaltedinc/framework/base"

const (
	tableStateIDConnectFightServer = base.StateID(hallStateIDCount + iota + 1)
	tableStateIDLoginFight
	tableStateIDTakeIn
	tableStateIDPlayerReady
)

type tableStateBase struct {
	table        *tableT
	user         *playerDataT
	stateManager *base.StateManager
}

func (state tableStateBase) OnEnter(prevStateID base.StateID) {

}

func (state tableStateBase) OnExit(nextStateID base.StateID) {

}

func (state tableStateBase) Update() {

}

func (state tableStateBase) OnMessage(data []byte) {

}

func (state tableStateBase) OnSuspend(nextStateID base.StateID) {

}

func (state tableStateBase) OnResume(prevStateID base.StateID) {

}
