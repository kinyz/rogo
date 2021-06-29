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
	CreatStateCluster(clusterId uint64, handle prehandle.Handle) error
	JoinStateCluster(clusterId uint64, handle prehandle.Handle) error

	CreatMemoryCluster(clusterId uint64, handle prehandle.Handle) error
	JoinMemoryCluster(clusterId uint64, handle prehandle.Handle) error

	CreatDiskCluster(clusterId uint64, handle prehandle.Handle) error
	JoinDiskCluster(clusterId uint64, handle prehandle.Handle) error

	SendMsg(clusterId uint64, msgId int64, data []byte) error
	SyncData(clusterId uint64, key string, value []byte) error
	GetSyncData(clusterId uint64, key string) ([]byte, error)
	RemoveSyncData(clusterId uint64, key string) error

	StopCluster(clusterId uint64) error

	AddNode(clusterId uint64, nodeId uint64, addr string) error

	GetSession(clusterId uint64) *client.Session

	RequestJoinCluster(grpcAddr string, clusterId uint64, key string) (*pb.ResponseJoinResult, error)
	StartOauthService(addr string, srv prehandle.RequestHandle) error

	GetNodeHost() *dragonboat.NodeHost
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

//提议数据
func (n *node) proposeData(clusterId uint64, proposeType pb.ProposeType, proposeId int64, data []byte) error {
	p := objPool.ProposePool.Get().(*pb.Propose)
	defer objPool.ProposePool.Put(p)
	p.TimesTamp = time.Now().UnixNano()
	p.Data = data
	p.NodeId = n.nh.ID()
	p.ProposeType = proposeType
	p.ProposeId = proposeId
	cs := n.nh.GetNoOPSession(clusterId)

	d, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = n.nh.SyncPropose(ctx, cs, d)
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

// SendMsg 发送消息
func (n *node) SendMsg(clusterId uint64, msgId int64, data []byte) error {
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

func (n *node) GetSyncData(clusterId uint64, key string) ([]byte, error) {
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
	//log.Println(read)
	//kvData:=objPool.KvPool.Get().(*pb.KvData)
	//err = proto.Unmarshal(read.([]byte), kvData)
	//if err != nil {
	//	return nil, err
	//}

	return read.([]byte), nil
}

func (n *node) RemoveSyncData(clusterId uint64, key string) error {
	return n.SyncData(clusterId, key, []byte(""))
}

func (n *node) CreatStateCluster(clusterId uint64, handle prehandle.Handle) error {
	initialMembers := make(map[uint64]string)
	initialMembers[1] = n.nh.RaftAddress()
	rc := config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}

	m := machine.NewStateMachine(handle)
	return n.nh.StartCluster(initialMembers, false, m.PreHandle, rc)

}
func (n *node) JoinStateCluster(clusterId uint64, handle prehandle.Handle) error {

	return n.nh.StartCluster(make(map[uint64]string), true, machine.NewStateMachine(handle).PreHandle, config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	})
}

func (n *node) CreatDiskCluster(clusterId uint64, handle prehandle.Handle) error {
	initialMembers := make(map[uint64]string)
	initialMembers[1] = n.nh.RaftAddress()
	rc := config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	m := machine.NewDiskKV(handle)
	if err := n.nh.StartOnDiskCluster(initialMembers, false, m.PreHandle, rc); err != nil {
		return err
	}
	return nil
}

func (n *node) JoinDiskCluster(clusterId uint64, handle prehandle.Handle) error {

	return n.nh.StartOnDiskCluster(make(map[uint64]string), true, machine.NewDiskKV(handle).PreHandle, config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	})
}

func (n *node) CreatMemoryCluster(clusterId uint64, handle prehandle.Handle) error {
	initialMembers := make(map[uint64]string)
	initialMembers[1] = n.nh.RaftAddress()
	rc := config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	m := machine.NewMemoryMachine(handle)
	if err := n.nh.StartCluster(initialMembers, false, m.PreStateHandle, rc); err != nil {
		return err
	}
	return nil
}

func (n *node) JoinMemoryCluster(clusterId uint64, handle prehandle.Handle) error {

	return n.nh.StartCluster(make(map[uint64]string), true, machine.NewMemoryMachine(handle).PreStateHandle, config.Config{
		NodeID:             n.nodeId,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	})
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

func (n *node) RequestJoinCluster(grpcAddr string, clusterId uint64, key string) (*pb.ResponseJoinResult, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewRequestClient(conn).RequestJoin(context.Background(), &pb.RequestJoinCluster{
		ClusterId:     clusterId,
		RaftAddresses: n.nh.RaftAddress(),
		Key:           key,
		NodeId:        n.nodeId,
	})

}

func (n *node) RequestJoin(ctx context.Context, req *pb.RequestJoinCluster) (*pb.ResponseJoinResult, error) {
	n.oauthLock.Lock()
	defer n.oauthLock.Unlock()
	err := n.reqDandle.Oauth(req)
	if err != nil {
		return nil, err
	}

	err = n.AddNode(req.GetClusterId(), req.GetNodeId(), req.GetRaftAddresses())
	if err != nil {
		return nil, err
	}
	return &pb.ResponseJoinResult{
		Status: 200,
		Msg:    "successful",
		//	NodeId: nodeId,
	}, nil

}
