package prehandle

import (
	"rogo/pb"
)

type Handle interface {
	SyncMessage(message *pb.Message)
}

type RequestHandle interface {
	Oauth(req *pb.RequestJoinCluster) (*pb.ResponseJoinResult, error)
}
