package machine

import (
	"errors"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"io"
	"rogo/objPool"
	"rogo/pb"
	"rogo/prehandle"
)

type MessageMachine struct {
	ClusterID uint64
	NodeID    uint64
	//Count     uint64
	handle prehandle.Handle
}

// NewMessageMachine creates and return a new ExampleStateMachine object.
func NewMessageMachine(handle prehandle.Handle) *MessageMachine {
	return &MessageMachine{
		handle: handle,
	}
}

func (s *MessageMachine) PreHandle(clusterID uint64,
	nodeID uint64) sm.IStateMachine {
	s.ClusterID = clusterID
	s.NodeID = nodeID
	return s
}

// Lookup performs local lookup on the ExampleStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *MessageMachine) Lookup(query interface{}) (interface{}, error) {
	//result := make([]byte, 8)
	//binary.LittleEndian.PutUint64(result, s.Count)

	//log.Println("我还是被执行了")

	return nil, errors.New("no storage")
}

// Update updates the object using the specified committed raft entry.
func (s *MessageMachine) Update(data []byte) (sm.Result, error) {

	p := objPool.ProposePool.Get().(*pb.Propose)
	switch p.GetProposeType() {
	case pb.ProposeType_SyncMessage:
		msg := objPool.MessagePool.Get().(*pb.Message)
		msg.ClusterId = s.ClusterID
		msg.NodeId = p.GetNodeId()
		msg.Data = p.GetData()
		msg.ProposeId = p.GetProposeId()
		s.handle.SyncMessage(msg)
		break
	}

	return sm.Result{Value: uint64(len(data))}, nil
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *MessageMachine) SaveSnapshot(w io.Writer,
	fc sm.ISnapshotFileCollection, done <-chan struct{}) error {
	// as shown above, the only state that can be saved is the Count variable
	// there is no external file in this IStateMachine example, we thus leave
	// the fc untouched
	//data := make([]byte, 8)
	//binary.LittleEndian.PutUint64(data, s.Count)
	//_, err := w.Write(data)
	return nil
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *MessageMachine) RecoverFromSnapshot(r io.Reader,
	files []sm.SnapshotFile,
	done <-chan struct{}) error {
	// restore the Count variable, that is the only state we maintain in this
	// example, the input files is expected to be empty
	//data, err := ioutil.ReadAll(r)
	//if err != nil {
	//	return err
	//}
	//v := binary.LittleEndian.Uint64(data)
	//s.Count = v
	return nil
}

// Close closes the IStateMachine instance. There is nothing for us to cleanup
// or release as this is a pure in memory data store. Note that the Close
// method is not guaranteed to be called as node can crash at any time.
func (s *MessageMachine) Close() error { return nil }
