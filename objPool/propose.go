package objPool

import (
	"rogo/pb"
	"sync"
)

var ProposePool =sync.Pool{
	New: func() interface{} {
		return &pb.Propose{
			ProposeType: 0,
			NodeId:      "",
			Data:        nil,
			ProposeId:   0,
			TimesTamp:   0,
		}
	},
}
