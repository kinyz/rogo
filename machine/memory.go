package machine

import (
	"errors"
	"github.com/golang/protobuf/proto"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"io"
	"log"
	"rogo/objPool"
	"rogo/pb"
	"sync"
)

func NewMemoryMachine(clusterID uint64,
	nodeID uint64) sm.IStateMachine {
	return &memoryMachine{ClusterID: clusterID, NodeID: nodeID}
}

type memoryMachine struct {
	//handle    prehandle.Handle
	ClusterID uint64
	NodeID    uint64
	kv        sync.Map
}

func (m *memoryMachine) Update(data []byte) (sm.Result, error) {
	p := objPool.ProposePool.Get().(*pb.Propose)

	err := proto.UnmarshalMerge(data, p)
	if err != nil {
		log.Println("解析错误:", err)
		return sm.Result{}, err
	}
	switch p.GetProposeType() {
	//case pb.ProposeType_SyncMessage:
	//	msg := objPool.MessagePool.Get().(*pb.Message)
	//	msg.ClusterId = m.ClusterID
	//	msg.NodeId = m.NodeID
	//	msg.ProposeId = p.GetProposeId()
	//	msg.Data = p.GetData()
	//	m.handle.SyncMessage(msg)
	//	break
	case pb.ProposeType_SyncData:
		kv := objPool.KvPool.Get().(*pb.KvData)
		err := proto.Unmarshal(p.GetData(), kv)
		if err != nil {
			log.Println(err)
			return sm.Result{Value: uint64(len(data))}, err
		}

		m.kv.Store(kv.GetKey(), kv.GetValue())
		break
	case pb.ProposeType_SyncRemoveData:
		kv := objPool.KvPool.Get().(*pb.KvData)
		err := proto.Unmarshal(p.GetData(), kv)
		if err != nil {
			return sm.Result{Value: uint64(len(data))}, err
		}
		m.kv.Delete(kv.GetKey())
		break
	}
	return sm.Result{Value: uint64(len(data))}, nil
}

func (m *memoryMachine) Lookup(query interface{}) (interface{}, error) {
	data, ok := m.kv.Load(string(query.([]byte)))
	if !ok {
		return nil, errors.New(string(query.([]byte)) + " is not existent")
	}
	return data.([]byte), nil

}

func (m *memoryMachine) SaveSnapshot(writer io.Writer, collection sm.ISnapshotFileCollection, i <-chan struct{}) error {

	return nil
}

func (m *memoryMachine) RecoverFromSnapshot(reader io.Reader, files []sm.SnapshotFile, i <-chan struct{}) error {

	return nil
}

func (m *memoryMachine) Close() error {

	return nil
}
