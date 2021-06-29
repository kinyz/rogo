package rogo

import (
	"rogo/node"
	"rogo/pb"
)

const (
	JoinRoleFollower = pb.JoinRole_Follower
	JoinRoleWitness  = pb.JoinRole_Witness
	JoinRoleObserver = pb.JoinRole_Observer

	StorageTypeDisk   = pb.StorageType_Disk
	StorageTypeMemory = pb.StorageType_Memory
)

type (
	Node               = node.Node
	Message            = *pb.Message
	RequestJoinCluster = *pb.RequestJoinCluster
)

func NewNode(addr, port string, nodeId uint64) Node {
	return node.NewNode(addr, port, nodeId)
}
