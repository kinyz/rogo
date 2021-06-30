package rogo

import (
	"rogo/node"
	"rogo/pb"
)

const (
	RoleCreator  = pb.Role_Creator
	RoleFollower = pb.Role_Follower
	RoleWitness  = pb.Role_Witness
	RoleObserver = pb.Role_Observer

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
