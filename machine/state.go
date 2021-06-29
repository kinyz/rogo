package machine

import (
	"errors"
	sm "github.com/lni/dragonboat/v3/statemachine"
	"io"
	"log"
	"rogo/objPool"
	"rogo/pb"
	"rogo/prehandle"
)

type StateMachine struct {
	ClusterID uint64
	NodeID    uint64
	//Count     uint64
	handle prehandle.Handle
}

// NewStateMachine creates and return a new ExampleStateMachine object.
func NewStateMachine(handle prehandle.Handle) *StateMachine {
	return &StateMachine{
		handle: handle,
	}
}

func (s *StateMachine) PreHandle(clusterID uint64,
	nodeID uint64) sm.IStateMachine {
	s.ClusterID = clusterID
	s.NodeID = nodeID
	return s
}

// Lookup performs local lookup on the ExampleStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *StateMachine) Lookup(query interface{}) (interface{}, error) {
	//result := make([]byte, 8)
	//binary.LittleEndian.PutUint64(result, s.Count)

	log.Println("我还是被执行了")

	return nil, errors.New("111")
}

// Update updates the object using the specified committed raft entry.
func (s *StateMachine) Update(data []byte) (sm.Result, error) {
	p := objPool.ProposePool.Get().(*pb.Propose)
	switch p.GetProposeType() {
	case pb.ProposeType_SyncMessage:
		s.handle.SyncMessage(&pb.Message{
			ClusterId: s.ClusterID,
			NodeId:    s.NodeID,
			ProposeId: p.GetProposeId(),
			Data:      p.GetData(),
		})
		break
	}

	return sm.Result{Value: uint64(len(data))}, nil
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *StateMachine) SaveSnapshot(w io.Writer,
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
func (s *StateMachine) RecoverFromSnapshot(r io.Reader,
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
func (s *StateMachine) Close() error { return nil }
