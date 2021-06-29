package rogo

import (
	"rogo/node"
	"rogo/pb"
)

type (
	Node               = node.Node
	Message            = *pb.Message
	RequestJoinCluster = *pb.RequestJoinCluster
)

func NewNode(addr, port string, nodeId uint64) Node {
	return node.NewNode(addr, port, nodeId)
}
