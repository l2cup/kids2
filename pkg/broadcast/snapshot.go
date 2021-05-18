package broadcast

import (
	"sync"
	"sync/atomic"

	"github.com/l2cup/kids2/internal/errors"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
)

type State struct {
	Token          string
	BitcakeBalance uint64
	Sent           map[uint64]Messages
	Recd           map[uint64]Messages
}

type Snapshots map[string]*Snapshot

type Snapshot struct {
	States      map[uint64]*State
	statesMutex sync.Mutex

	Waiting uint64
	Got     uint64
}

func (s *Snapshot) AddState(st *State, ID uint64) {
	defer s.statesMutex.Unlock()
	s.statesMutex.Lock()

	atomic.AddUint64(&s.Got, 1)
	s.States[ID] = st
}

func (s *Snapshot) GetStates() (map[uint64]*State, errors.Error) {
	if !s.Finished() {
		return nil, errors.New(
			"snapshot not done",
			"total", atomic.LoadUint64(&s.Waiting),
			"got", atomic.LoadUint64(&s.Got),
		)
	}
	return s.States, errors.Nil()
}

func (s *Snapshot) Finished() bool {
	return s.Waiting == s.Got
}

func (s *State) Proto() *nodepb.State {
	return &nodepb.State{
		Token:    s.Token,
		Bitcakes: s.BitcakeBalance,
		Sent:     MsgMapToKVStatePairs(s.Sent),
		Received: MsgMapToKVStatePairs(s.Recd),
	}
}

func SnapshotStateFromProto(s *nodepb.State) *State {
	return &State{
		Token:          s.GetToken(),
		BitcakeBalance: s.GetBitcakes(),
		Sent:           KVStatePairsToMap(s.GetSent()...),
		Recd:           KVStatePairsToMap(s.GetReceived()...),
	}
}

func KVStatePairsToMap(skvs ...*nodepb.StateKVPairs) map[uint64]Messages {
	msgMap := make(map[uint64]Messages, len(skvs))
	for _, skv := range skvs {
		msgMap[skv.GetKey()] = MessagesFromProto(skv.GetValues()...)
	}
	return msgMap
}

func MsgMapToKVStatePairs(msgMap map[uint64]Messages) []*nodepb.StateKVPairs {
	kvPairs := make([]*nodepb.StateKVPairs, 0, len(msgMap))
	for k, v := range msgMap {
		kvPair := &nodepb.StateKVPairs{
			Key:    k,
			Values: v.Proto(),
		}
		kvPairs = append(kvPairs, kvPair)
	}

	return kvPairs
}

type SnapshotRequest struct {
	Token string
}

func (r *SnapshotRequest) Proto() *nodepb.SnapshotRequest {
	return &nodepb.SnapshotRequest{
		Token: r.Token,
	}
}

func RequestFromProto(req *nodepb.SnapshotRequest) *SnapshotRequest {
	return &SnapshotRequest{
		Token: req.Token,
	}
}
