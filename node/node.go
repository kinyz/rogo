package node

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/lni/dragonboat/v3"
	"github.com/lni/dragonboat/v3/client"
	"github.com/lni/dragonboat/v3/config"
	"path/filepath"
	"rogo/machine"
	"rogo/objPool"
	"rogo/pb"
	"rogo/prehandle"
	"rogo/storage"
	"time"
)

type Node interface {

	StartStateCluster(clusterId uint64,message prehandle.Message)error
	StartKvCluster(clusterId uint64,message prehandle.Message)error

	JoinStateCluster(clusterId uint64,message prehandle.Message)error
	JoinKvCluster(clusterId uint64,message prehandle.Message)error

	SendMsg(clusterId uint64,msgId int64 ,data []byte)error
	SyncData(clusterId uint64,key string,value []byte)error
	GetSyncData(clusterId uint64,key string)([]byte,error)
	RemoveSyncData(clusterId uint64,key string)error

	LogoutCluster(clusterId uint64)error

	AddNode(clusterId uint64,nodeId uint64,addr string)error

	GetSession(clusterId uint64)*client.Session
}

func NewNode(addr,port string,nodeId uint64)Node{
	datadir := filepath.Join(
		"sync-data",
		"key-data",
		fmt.Sprintf("node%d", nodeId))

	nhc := config.NodeHostConfig{
		WALDir:         datadir,
		NodeHostDir:    datadir,
		RTTMillisecond: 200,
		RaftAddress:    addr+":"+port,
		//	AddressByNodeHostID:true,
		ListenAddress: "0.0.0.0:" + port,
	}
	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		panic(err)
	}
	return &node{nh: nh,nodeId: nodeId}
}
type node struct {
	nh	*dragonboat.NodeHost
	nodeId uint64
}


//提议数据
func (n*node)proposeData(clusterId uint64,proposeType pb.ProposeType,proposeId int64,data []byte)error{
	p:=objPool.ProposePool.Get().(*pb.Propose)
	defer objPool.ProposePool.Put(p)
	p.TimesTamp=time.Now().UnixNano()
	p.Data = data
	p.NodeId = n.nh.ID()
	p.ProposeType = proposeType
	p.ProposeId = proposeId
	cs := n.nh.GetNoOPSession(clusterId)

	d,err:=proto.Marshal(p)
	if err!=nil{
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = n.nh.SyncPropose(ctx,cs,d)
	if err != nil {
		cancel()
		return err
	}
	cancel()
	return nil
}

// SendMsg 发送消息
func (n*node)SendMsg(clusterId uint64,msgId int64 ,data []byte)error{
	return n.proposeData(clusterId,pb.ProposeType_SyncMessage,msgId,data)
}

func (n*node)SyncData(clusterId uint64,key string,value []byte)error{
	kvData:=objPool.KvPool.Get().(*pb.KvData)
	defer objPool.KvPool.Put(kvData)
	kvData.Key = key
	kvData.Value = value
	data,err:=proto.Marshal(kvData)
	if err!=nil{
		return err
	}
	return n.proposeData(clusterId,pb.ProposeType_SyncData,0,data)
}

func (n*node)GetSyncData(clusterId uint64,key string)([]byte,error){
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	read, err := n.nh.SyncRead(ctx, clusterId, []byte(key))
	if err != nil {
		cancel()
		return nil, err
	}
	cancel()
	if string(read.([]byte))=="" {
		return nil,errors.New("key "+key+" is not exist")
	}
	//log.Println(read)
	//kvData:=objPool.KvPool.Get().(*pb.KvData)
	//err = proto.Unmarshal(read.([]byte), kvData)
	//if err != nil {
	//	return nil, err
	//}

	return read.([]byte),nil
}

func (n*node)RemoveSyncData(clusterId uint64,key string)error{
	return n.SyncData(clusterId,key,[]byte(""))
}

func (n*node)StartStateCluster(clusterId uint64,message prehandle.Message)error{
	initialMembers := make(map[uint64]string)
	initialMembers[1] = n.nh.RaftAddress()
	rc := config.Config{
		NodeID:             n.nodeId ,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}

	m:=machine.NewStateMachine(message)
	return n.nh.StartCluster(initialMembers, false,m.PreHandle , rc)

}

func (n*node)StartKvCluster(clusterId uint64,message prehandle.Message)error{
	initialMembers := make(map[uint64]string)
	initialMembers[1] = n.nh.RaftAddress()
	rc := config.Config{
		NodeID:             n.nodeId ,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	m:=storage.NewDiskKV(message)
	if err := n.nh.StartOnDiskCluster(initialMembers, false,m.PreHandle, rc); err != nil {
		return err
	}
	return nil
}

func (n*node)JoinStateCluster(clusterId uint64,message prehandle.Message)error{
	initialMembers := make(map[uint64]string)
	rc := config.Config{
		NodeID:             n.nodeId ,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	m:=machine.NewStateMachine(message)
	return n.nh.StartCluster(initialMembers, true,m.PreHandle , rc)


}

func (n*node)JoinKvCluster(clusterId uint64,message prehandle.Message)error{
	initialMembers := make(map[uint64]string)
	rc := config.Config{
		NodeID:             n.nodeId ,
		ClusterID:          clusterId,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    10,
		CompactionOverhead: 5,
	}
	m:=storage.NewDiskKV(message)
	if err := n.nh.StartOnDiskCluster(initialMembers, true,m.PreHandle, rc); err != nil {
		return err
	}

	return nil
}


func (n*node)LogoutCluster(clusterId uint64)error{
	return n.nh.StopCluster(clusterId)
}

func (n*node)AddNode(clusterId uint64,nodeId uint64,addr string)error{
	_, err := n.nh.RequestAddNode(clusterId,nodeId,addr,0, 5*time.Second)
	if err != nil {
		return err
	}
	return nil
}


func (n*node)GetSession(clusterId uint64)*client.Session{
	return n.nh.GetNoOPSession(clusterId)
}