package broadcast

import (
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/vc"
)

type Type string

const (
	TypeTransaction Type = "TRASACTION"
	TypeSnapshot    Type = "SNAPSHOT"
)

type Message struct {
	ID     uint64
	From   uint64
	To     uint64
	VClock *vc.VectorClock
	Type   Type
	Data   interface{}
}

func (m *Message) Proto() *nodepb.Message {
	&nodepb.Message{}
}

type Transaction struct {
	Bitcakes uint64
}

func (t *Transaction) Proto() *nodepb.Transaction {
	return &nodepb.Transaction{
		Bitcakes: t.Bitcakes,
	}
}

func TransactionFromProto(t *nodepb.Transaction) *Transaction {
	return &Transaction{
		Bitcakes: t.GetBitcakes(),
	}
}

type Snapshot struct {
	Sent     []uint64
	Received []uint64
}

func (s *Snapshot) Proto() *nodepb.Snapshot {
	return &nodepb.Snapshot{
		Sent:     s.Sent,
		Received: s.Received,
	}
}

func SnapshotFromProto(s *nodepb.Snapshot) *Snapshot {
	return &Snapshot{
		Sent:     s.GetSent(),
		Received: s.GetReceived(),
	}
}
