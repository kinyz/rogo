package objPool

import (
	"rogo/pb"
	"sync"
)

var MessagePool = sync.Pool{
	New: func() interface{} {
		return &pb.Message{}
	},
}
