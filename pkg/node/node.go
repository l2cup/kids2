package node

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/l2cup/kids2/internal/errors"
	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/broadcast"
	"github.com/l2cup/kids2/pkg/network"
	"github.com/l2cup/kids2/pkg/vc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Node struct {
	nodepb.UnimplementedNodeServer
	logger *log.Logger

	broadcast *broadcast.Broadcast

	bitcakeBalance uint64
	//bitcakeMutex   sync.Mutex

	vclock *vc.VectorClock

	sent      map[uint64][]*broadcast.Message
	sentMutex sync.Mutex

	recd      map[uint64][]*broadcast.Message
	recdMutex sync.Mutex

	commited uint64

	bufferMutex sync.Mutex
	buffer      []*broadcast.Message

	bootstrapInfo  *network.Info
	networkInfo    *network.Info
	connectedNodes []*network.Info
	totalNodeCount uint64
}

func (n *Node) BroadcastTransaction(ctx context.Context, to uint64, change uint64) errors.Error {
	connected := false
	for _, node := range n.connectedNodes {
		if node.ID == to {
			connected = true
			break
		}
	}

	if !connected {
		return errors.New(
			"no connection to node",
			errors.InternalServerError,
			"node_id", to,
		)
	}

	for {
		b := atomic.LoadUint64(&n.bitcakeBalance)
		if b < change {
			return errors.New("not enough bitcakes", errors.InternalServerError)
		}
		newBalance := b - change
		if swapped := atomic.CompareAndSwapUint64(&n.bitcakeBalance, b, newBalance); swapped {
			n.logger.Debug("bitcake swap",
				"node", n.networkInfo.ID,
				"to_node", to,
				"old_balance", b,
				"new_balance", newBalance,
				"change", change,
			)
			break
		}
		n.logger.Debug("bitcake swap for", "bc_balance", b, "new_balance", newBalance)
	}

	vc := n.vclock.Copy()
	time, _ := vc.TimeUint64(n.networkInfo.ID)

	msg := &broadcast.Message{
		ID:     time,
		From:   n.networkInfo.ID,
		To:     to,
		VClock: vc,
		Type:   broadcast.TypeTransaction,
		Data: &broadcast.Transaction{
			Bitcakes: change,
		},
	}

	n.logger.Debug("broadcasting message to nodes", "msg", msg, "nodes", n.connectedNodes)
	n.broadcast.Broadcast(msg, n.connectedNodes)

	defer n.sentMutex.Unlock()
	n.sentMutex.Lock()
	n.sent[n.networkInfo.ID] = append(n.sent[n.networkInfo.ID], msg)
	n.vclock.TickUint64(n.networkInfo.ID)

	return errors.Nil()
}

func (n *Node) Transaction(ctx context.Context, msgpb *nodepb.Message) (*emptypb.Empty, error) {
	time.Sleep(time.Duration((rand.Intn(5000) + 50)) * time.Millisecond)

	if msgpb.From == n.networkInfo.ID {
		return &emptypb.Empty{}, nil
	}

	defer n.bufferMutex.Unlock()
	n.bufferMutex.Lock()

	msg := broadcast.MessageFromProto(msgpb)
	n.buffer = append(n.buffer, msg)
	n.commit()

	return &emptypb.Empty{}, nil
}

func (n *Node) commit() {
	defer n.recdMutex.Unlock()
	n.recdMutex.Lock()

	for i, msg := range n.buffer {
		if msg == nil {
			continue
		}

		if !n.vclock.IsDescendantOrEqual(msg.VClock) {
			n.logger.Info("not descendant or equal", "vc", n.vclock.Map(), "msg vc", msg.VClock.Map(), "node", n.networkInfo.ID)
			continue
		}

		if n.alreadyRecd(msg) {
			n.commited += 1
			n.buffer[i] = nil
			continue
		}

		switch msg.Type {
		case broadcast.TypeTransaction:
			n.commitTx(msg)
		case broadcast.TypeSnapshot:
		}

		n.logger.Debug("rebroadcasting message", "msg", msg, "node", n.networkInfo.ID)
		n.broadcast.Broadcast(msg, n.connectedNodes)
	}

	//if n.commited >= 5 {
	//n.commited = 0
	//newBuffer := make([]*broadcast.Message, 0, len(n.buffer)-int(n.commited))
	//for _, msg := range n.buffer {
	//if msg == nil {
	//continue
	//}
	//newBuffer = append(newBuffer, msg)
	//}
	//n.buffer = newBuffer
	//}
}

func (n *Node) alreadyRecd(msg *broadcast.Message) bool {
	for _, recdMsg := range n.recd[msg.From] {
		if recdMsg.ID == msg.ID {
			return true
		}
	}

	return false
}

func (n *Node) commitTx(msg *broadcast.Message) {
	tx, ok := msg.Data.(*broadcast.Transaction)
	if !ok {
		n.logger.Error(
			"couldn't commit transaction, data type not transaction",
			"type", fmt.Sprintf("%T", msg.Data),
		)
	}

	if msg.To == n.networkInfo.ID {
		oldBalance := atomic.LoadUint64(&n.bitcakeBalance)
		atomic.AddUint64(&n.bitcakeBalance, tx.Bitcakes)
		n.logger.Debug(
			"received transaction",
			"node", n.networkInfo.ID,
			"old_balance", oldBalance,
			"new_balance", atomic.LoadUint64(&n.bitcakeBalance),
			"change", tx.Bitcakes,
		)
	}

	n.recd[msg.From] = append(n.recd[msg.From], msg)
	n.vclock.TickUint64(msg.From)
	n.logger.Debug(
		"commited transaction",
		"msg", msg,
		"msg_vc", msg.VClock.Map(),
		"node", n.networkInfo.ID,
		"vc", n.vclock.Map(),
	)
}
