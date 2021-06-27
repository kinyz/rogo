package prehandle

import "rogo/pb"

type Message interface {
	SyncMessage(message *pb.Message)
}
