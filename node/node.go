package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/config"
	"google.golang.org/grpc"
	"log"
	"net"
	"path/filepath"
	"rogo/machine"
	"rogo/objPool"
	"rogo/pb"
	"rogo/prehandle"
	"sync"
	"time"
)

type Node interface {

	// StartStorageCluster 创建存储集群
	StartStorageCluster(clusterId uint64, machineType pb.StorageType, role pb.Role) error

	// StartMessageCluster 创建消息集群
	StartMessageCluster(clusterId uint64, role pb.Role, handle prehandle.Handle) error

	// SendMessage 发送消息 仅消息集群可用
	SendMessage(clusterId uint64, msgId int64, data []byte) error

	SyncData(clusterId uint64, key string, value []byte) error
	GetData(clusterId uint64, key string) ([]byte, error)
	RemoveData(clusterId uint64, key string) error

	StopCluster(clusterId uint64) error

	AddNode(clusterId uint64, nodeId uint64, addr string) error
	AddWitnessNode(clusterId uint64, nodeId uint64, addr string) error
	AddObserverNode(clusterId uint64, nodeId uint64, addr string) error

	// BanNode 删除Node 此Node进入黑名单
	BanNode(clusterId uint64, node uint64) error

	RequestJoinCluster(grpcAddr string, clusterId uint64, key string, role pb.Role) (*pb.ResponseJoinResult, error)

	StartOauthService(addr string, srv prehandle.RequestHandle) error

	GetRaftId() string
}

func NewNode(addr, port string, nodeId uint64) Node {
	datadir := filepath.Join(
		"sync-data",
		"key-data",
		fmt.Sprintf("node%d", nodeId))

	nhc := config.NodeHostConfig{
		WALDir:         datadir,
		NodeHostDir:    datadir,
		RTTMillisecond: 200,
		RaftAddress:    addr + ":" + port,
		//	AddressByNodeHostID:true,
		ListenAddress: "0.0.0.0:" + port,
	}
	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		panic(err)
	}

	//nh.StartOnDiskCluster(nil,false,machine.Machines{},nil)
	return &node{nh: nh, nodeId: nodeId}
}

type node struct {
	nh     *dragonboat.NodeHost
	nodeId uint64

	oauthLock sync.Mutex

	reqDandle prehandle.RequestHandle
}

func (n *node) StartStorageCluster(clusterId uint64, machineType pb.StorageType, role pb.Role) error {

	join := false
	initialMembers := make(map[uint64]string)
	rc := config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	switch role {
	case pb.Role_Creator:
		initialMembers[1] = n.nh.RaftAddress()
		break
	case pb.Role_Follower:
		join = true
		break
	case pb.Role_Witness:
		rc.SnapshotEntries = 0
		rc.IsWitness = true
		join = true
	case pb.Role_Observer:
		rc.IsObserver = true
		join = true
	default:
		return errors.New("role is error")
	}

	switch machineType {
	case pb.StorageType_Disk:
		return n.nh.StartOnDiskCluster(initialMembers, join, machine.NewDiskKV, rc)
	case pb.StorageType_Memory:
		return n.nh.StartCluster(initialMembers, join, machine.NewMemoryMachine, rc)
	}
	return errors.New("storage type is error")

}

func (n *node) StartMessageCluster(clusterId uint64, role pb.Role, handle prehandle.Handle) error {
	initialMembers := make(map[uint64]string)

	rc := config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	join := false
	switch role {
	case pb.Role_Creator:
		initialMembers[1] = n.nh.RaftAddress()
		break
	case pb.Role_Follower:
		join = true
		break
	case pb.Role_Witness:
		rc.SnapshotEntries = 0
		rc.IsWitness = true
		join = true
	case pb.Role_Observer:
		rc.IsObserver = true
		join = true
	default:
		return errors.New("role is error")
	}
	m := machine.NewMessageMachine(handle)
	return n.nh.StartCluster(initialMembers, join, m.PreHandle, rc)

}

//提议数据
func (n *node) proposeData(clusterId uint64, proposeType pb.ProposeType, proposeId int64, data []byte) error {
	p := objPool.ProposePool.Get().(*pb.Propose)
	defer objPool.ProposePool.Put(p)
	p.TimesTamp = time.Now().UnixNano()
	p.Data = data
	p.NodeId = n.GetRaftId()
	p.ProposeType = proposeType
	p.ProposeId = proposeId
	d, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = n.nh.SyncPropose(ctx, n.nh.GetNoOPSession(clusterId), d)
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}

func (n *node) GetNodeHost() *dragonboat.NodeHost {
	return n.nh
}

// SendMessage 发送消息
func (n *node) SendMessage(clusterId uint64, msgId int64, data []byte) error {
	return n.proposeData(clusterId, pb.ProposeType_SyncMessage, msgId, data)
}

func (n *node) SyncData(clusterId uint64, key string, value []byte) error {
	kvData := objPool.KvPool.Get().(*pb.KvData)
	defer objPool.KvPool.Put(kvData)
	kvData.Key = key
	kvData.Value = value
	data, err := proto.Marshal(kvData)
	if err != nil {
		return err
	}

	return n.proposeData(clusterId, pb.ProposeType_SyncData, 0, data)
}

func (n *node) GetData(clusterId uint64, key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	read, err := n.nh.SyncRead(ctx, clusterId, []byte(key))
	if err != nil {
		cancel()
		return nil, err
	}
	cancel()
	if string(read.([]byte)) == "" {
		return nil, errors.New("key " + key + " is not exist")
	}

	return read.([]byte), nil
}

func (n *node) RemoveData(clusterId uint64, key string) error {
	kvData := objPool.KvPool.Get().(*pb.KvData)
	defer objPool.KvPool.Put(kvData)
	kvData.Key = key
	kvData.Value = []byte("")
	data, err := proto.Marshal(kvData)
	if err != nil {
		return err
	}
	return n.proposeData(clusterId, pb.ProposeType_SyncRemoveData, 0, data)
}

func (n *node) StopCluster(clusterId uint64) error {
	return n.nh.StopCluster(clusterId)
}

func (n *node) AddNode(clusterId uint64, nodeId uint64, addr string) error {
	_, err := n.nh.RequestAddNode(clusterId, nodeId, addr, 0, 5*time.Second)
	if err != nil {
		return err
	}

	return nil
}
func (n *node) AddWitnessNode(clusterId uint64, nodeId uint64, addr string) error {
	_, err := n.nh.RequestAddWitness(clusterId, nodeId, addr, 0, 5*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (n *node) AddObserverNode(clusterId uint64, nodeId uint64, addr string) error {
	_, err := n.nh.RequestAddObserver(clusterId, nodeId, addr, 0, 5*time.Second)
	if err != nil {
		return err
	}

	return nil
}

func (n *node) GetRaftId() string {
	return n.nh.ID()
}

func (n *node) GetSession(clusterId uint64) *client.Session {
	return n.nh.GetNoOPSession(clusterId)
}

func (n *node) StartOauthService(addr string, srv prehandle.RequestHandle) error {
	if n.reqDandle != nil {
		return errors.New("do not turn on repeatedly ！")
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterRequestServer(s, n)
	go func() {
		err = s.Serve(ln)
		if err != nil {
			fmt.Println("用户端口出错", err)
		}
	}()
	log.Println("oauth service start: ", addr)

	n.reqDandle = srv
	return nil
}

func (n *node) RequestJoinCluster(grpcAddr string, clusterId uint64, key string, role pb.Role) (*pb.ResponseJoinResult, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewRequestClient(conn).RequestJoin(context.Background(), &pb.RequestJoinCluster{
		ClusterId:     clusterId,
		RaftAddresses: n.nh.RaftAddress(),
		Key:           key,
		NodeId:        n.nodeId,
		JoinRole:      role,
	})
}

func (n *node) RequestJoin(ctx context.Context, req *pb.RequestJoinCluster) (*pb.ResponseJoinResult, error) {
	n.oauthLock.Lock()
	defer n.oauthLock.Unlock()
	resp, err := n.reqDandle.Oauth(req)
	if err != nil {
		return nil, err
	}

	switch req.GetJoinRole() {
	case pb.Role_Follower:
		err = n.AddNode(resp.GetClusterId(), req.GetNodeId(), req.GetRaftAddresses())
		if err != nil {
			return nil, err
		}
		break
	case pb.Role_Witness:
		err = n.AddWitnessNode(resp.GetClusterId(), req.GetNodeId(), req.GetRaftAddresses())
		if err != nil {
			return nil, err
		}
		break
	case pb.Role_Observer:
		err = n.AddObserverNode(resp.GetClusterId(), req.GetNodeId(), req.GetRaftAddresses())
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("role error")
	}

	return resp, nil

}

func (n *node) BanNode(clusterId uint64, node uint64) error {
	_, err := n.nh.RequestDeleteNode(clusterId, node, 0, 5*time.Second)
	if err != nil {
		return err
	}
	return nil
}
