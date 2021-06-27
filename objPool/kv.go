package objPool

import (
	"rogo/pb"
	"sync"
)

var KvPool =sync.Pool{
	New: func() interface{} {
		return &pb.KvData{
			}
	},
}