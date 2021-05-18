package broadcast

import (
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/vc"
)

type Type string

const (
	TypeTransaction     Type = "TRASACTION"
	TypeSnapshotRequest Type = "SNAPSHOT_REQUEST"
	TypeSnapshotState   Type = "SNAPSHOT_STATE"
)

type Message struct {
	ID     uint64
	From   uint64
	To     uint64
	VClock *vc.VectorClock

	Type Type
	Data interface{}
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
	case TypeSnapshotState:
		pbmsg.Type = m.protoSnapshotState()
	case TypeSnapshotRequest:
		pbmsg.Type = m.protoSnapshotRequest()
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

func (m *Message) protoSnapshotState() *nodepb.Message_SnapshotState {
	ss, ok := m.Data.(*State)
	if !ok {
		return nil
	}

	return &nodepb.Message_SnapshotState{
		SnapshotState: ss.Proto(),
	}
}

func (m *Message) protoSnapshotRequest() *nodepb.Message_SnapshotRequest {
	sr, ok := m.Data.(*SnapshotRequest)
	if !ok {
		return nil
	}

	return &nodepb.Message_SnapshotRequest{
		SnapshotRequest: sr.Proto(),
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
	if x, ok := t.(*nodepb.Message_SnapshotState); ok {
		msg.Type = TypeSnapshotState
		msg.Data = SnapshotStateFromProto(x.SnapshotState)
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

type Messages []*Message

func (mm Messages) Proto() []*nodepb.Message {
	mmpb := make([]*nodepb.Message, 0, len(mm))
	for _, m := range mm {
		mmpb = append(mmpb, m.Proto())
	}

	return mmpb
}

func MessagesFromProto(mmpb ...*nodepb.Message) Messages {
	mm := make(Messages, 0, len(mmpb))
	for _, m := range mmpb {
		mm = append(mm, MessageFromProto(m))
	}

	return mm
}

func (mm Messages) Contains(msg *Message) bool {
	for _, m := range mm {
		if m.ID == msg.ID {
			return true
		}
	}
	return false
}

func (mm Messages) Append(msg *Message) Messages {
	mm = append(mm, msg)
	return mm
}
