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
	pbmsg := &nodepb.Message{
		Id:     m.ID,
		From:   m.From,
		To:     m.To,
		Vclock: m.VClock.Proto(),
	}

	switch m.Type {
	case TypeTransaction:
		pbmsg.Type = m.protoTransaction()
	case TypeSnapshot:
		pbmsg.Type = m.protoSnapshot()
	}
	return pbmsg
}

func (m *Message) protoTransaction() *nodepb.Message_Transaction {
	tx, ok := m.Data.(*Transaction)
	if !ok {
		return nil
	}

	return &nodepb.Message_Transaction{
		Transaction: tx.Proto(),
	}
}

func (m *Message) protoSnapshot() *nodepb.Message_Snapshot {
	ss, ok := m.Data.(*Snapshot)
	if !ok {
		return nil
	}

	return &nodepb.Message_Snapshot{
		Snapshot: ss.Proto(),
	}
}

func MessageFromProto(msgpb *nodepb.Message) *Message {
	msg := &Message{
		ID:     msgpb.GetId(),
		From:   msgpb.GetFrom(),
		To:     msgpb.GetTo(),
		VClock: vc.VectorClockFromProto(msgpb.GetVclock()),
	}

	t := msgpb.GetType()
	if x, ok := t.(*nodepb.Message_Transaction); ok {
		msg.Type = TypeTransaction
		msg.Data = TransactionFromProto(x.Transaction)
		return msg
	}
	if x, ok := t.(*nodepb.Message_Snapshot); ok {
		msg.Type = TypeSnapshot
		msg.Data = SnapshotFromProto(x.Snapshot)
		return msg
	}

	return msg
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
	Bitcakes uint64
	Sent     []uint64
	Received []uint64
}

func (s *Snapshot) Proto() *nodepb.Snapshot {
	return &nodepb.Snapshot{
		Bitcakes: s.Bitcakes,
		Sent:     s.Sent,
		Received: s.Received,
	}
}

func SnapshotFromProto(s *nodepb.Snapshot) *Snapshot {
	return &Snapshot{
		Bitcakes: s.GetBitcakes(),
		Sent:     s.GetSent(),
		Received: s.GetReceived(),
	}
}
