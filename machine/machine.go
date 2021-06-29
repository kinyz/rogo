package machine

import sm "github.com/lni/dragonboat/v3/statemachine"

type Machine interface {
	PreStateHandle(clusterID uint64,
		nodeID uint64) sm.IStateMachine
}
