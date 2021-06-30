package rogo

import (
	"rogo/node"
	"rogo/pb"
)

const (
	// RoleCreator 集群创建者
	RoleCreator = pb.Role_Creator
	// RoleFollower 集群跟随者
	RoleFollower = pb.Role_Follower
	// RoleWitness 集群见证者
	RoleWitness = pb.Role_Witness
	// RoleObserver 集群观察者
	RoleObserver = pb.Role_Observer

	// StorageTypeDisk 硬盘存储
	StorageTypeDisk = pb.StorageType_Disk
	// StorageTypeMemory 内存存储
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
