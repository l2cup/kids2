package node

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
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

	sent      map[uint64]broadcast.Messages
	sentMutex sync.Mutex

	recd      map[uint64]broadcast.Messages
	recdMutex sync.Mutex

	processed      map[uint64]broadcast.Messages
	processedMutex sync.Mutex

	commited uint64

	bufferMutex sync.Mutex
	buffer      []*broadcast.Message

	snapshots     broadcast.Snapshots
	snapshotMutex sync.Mutex

	broadcastLock sync.Mutex

	bootstrapInfo  *network.Info
	networkInfo    *network.Info
	connectedNodes []*network.Info
	totalNodeCount uint64
}

func (n *Node) State() *broadcast.State {
	defer n.recdMutex.Unlock()
	defer n.sentMutex.Unlock()
	n.recdMutex.Lock()
	n.sentMutex.Lock()

	sent := make(map[uint64]broadcast.Messages, len(n.sent))
	recd := make(map[uint64]broadcast.Messages, len(n.recd))
	for k, v := range n.sent {
		sent[k] = append(sent[k], v...)
	}

	for k, v := range n.recd {
		recd[k] = append(recd[k], v...)
	}

	bitcakes := atomic.LoadUint64(&n.bitcakeBalance)

	return &broadcast.State{
		BitcakeBalance: bitcakes,
		Sent:           sent,
		Recd:           recd,
	}
}

func (n *Node) BroadcastSnapshotRequest(ctx context.Context) errors.Error {
	defer n.snapshotMutex.Unlock()
	n.snapshotMutex.Lock()

	defer n.broadcastLock.Unlock()
	n.broadcastLock.Lock()

	for _, v := range n.snapshots {
		if !v.Finished() {
			return errors.New(
				"snapshot in progress",
				errors.BadRequestError,
			)
		}
	}

	token := uuid.New().String()
	n.snapshots[token] = &broadcast.Snapshot{
		States:  make(map[uint64]*broadcast.State, n.totalNodeCount),
		Waiting: n.totalNodeCount,
	}

	state := n.State()
	n.snapshots[token].AddState(state, n.networkInfo.ID)

	vc := n.vclock.Copy()
	time, _ := vc.TimeUint64(n.networkInfo.ID)

	msg := &broadcast.Message{
		ID:     time,
		From:   n.networkInfo.ID,
		To:     0xdead,
		VClock: vc,
		Type:   broadcast.TypeSnapshotRequest,
		Data: &broadcast.SnapshotRequest{
			Token: token,
		},
	}

	defer n.processedMutex.Unlock()
	n.processedMutex.Lock()

	n.processed[n.networkInfo.ID] = append(n.processed[n.networkInfo.ID], msg)

	n.broadcast.Broadcast(msg, n.connectedNodes)
	return errors.Nil()
}

func (n *Node) BroadcastTransaction(ctx context.Context, to uint64, change uint64) errors.Error {
	defer n.broadcastLock.Unlock()
	n.broadcastLock.Lock()

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

	n.sent[msg.To] = append(n.sent[msg.To], msg)

	defer n.processedMutex.Unlock()
	n.processedMutex.Lock()

	n.processed[n.networkInfo.ID] = append(n.processed[n.networkInfo.ID], msg)

	n.vclock.TickUint64(n.networkInfo.ID)

	return errors.Nil()
}

func (n *Node) SendSnapshot(ctx context.Context, msgpb *nodepb.Message) (*emptypb.Empty, error) {
	defer n.bufferMutex.Unlock()
	n.bufferMutex.Lock()

	msg := broadcast.MessageFromProto(msgpb)
	n.buffer = append(n.buffer, msg)
	n.commit()

	return &emptypb.Empty{}, nil
}

func (n *Node) RequestSnapshot(ctx context.Context, msgpb *nodepb.Message) (*emptypb.Empty, error) {
	time.Sleep(time.Duration((rand.Intn(2500) + 50)) * time.Millisecond)

	defer n.bufferMutex.Unlock()
	n.bufferMutex.Lock()

	msg := broadcast.MessageFromProto(msgpb)
	n.buffer = append(n.buffer, msg)
	n.commit()

	return &emptypb.Empty{}, nil
}

func (n *Node) Transaction(ctx context.Context, msgpb *nodepb.Message) (*emptypb.Empty, error) {
	time.Sleep(time.Duration((rand.Intn(5000) + 50)) * time.Millisecond)

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

	defer n.processedMutex.Unlock()
	n.processedMutex.Lock()

	for i, msg := range n.buffer {
		if msg == nil {
			continue
		}

		if !n.vclock.IsDescendantOrEqual(msg.VClock) {
			n.logger.Info("not descendant or equal", "vc", n.vclock.Map(), "msg vc", msg.VClock.Map(), "node", n.networkInfo.ID)
			continue
		}

		if n.alreadyProcessed(msg) {
			n.commited += 1
			n.buffer[i] = nil
			continue
		}

		switch msg.Type {
		case broadcast.TypeTransaction:
			n.commitTx(msg)
		case broadcast.TypeSnapshotRequest:
			n.commitSnapshotRequest(msg)
		case broadcast.TypeSnapshotState:
			n.commitSnapshotState(msg)
		}

		n.logger.Debug("rebroadcasting message", "msg", msg, "node", n.networkInfo.ID)
		n.broadcast.Broadcast(msg, n.connectedNodes)
	}

	if n.commited >= 5 {
		n.commited = 0
		newBuffer := make([]*broadcast.Message, 0, len(n.buffer)-int(n.commited))
		for _, msg := range n.buffer {
			if msg == nil {
				continue
			}
			newBuffer = append(newBuffer, msg)
		}
		n.buffer = newBuffer
	}
}

func (n *Node) alreadyProcessed(msg *broadcast.Message) bool {
	return n.processed[msg.From].Contains(msg)
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
		n.recd[msg.From] = append(n.recd[msg.From], msg)
		n.logger.Debug(
			"received transaction",
			"node", n.networkInfo.ID,
			"old_balance", oldBalance,
			"new_balance", atomic.LoadUint64(&n.bitcakeBalance),
			"change", tx.Bitcakes,
		)
	}

	n.processed[msg.From] = append(n.processed[msg.From], msg)
	n.vclock.TickUint64(msg.From)
	n.logger.Debug(
		"commited transaction",
		"msg", msg,
		"msg_vc", msg.VClock.Map(),
		"node", n.networkInfo.ID,
		"vc", n.vclock.Map(),
	)
}

func (n *Node) commitSnapshotRequest(msg *broadcast.Message) {
	req, ok := msg.Data.(*broadcast.SnapshotRequest)
	if !ok {
		n.logger.Error(
			"couldn't commit snapshot request, data type not snapshot request",
			"type", fmt.Sprintf("%T", msg.Data),
		)
	}

	n.processed[msg.From] = append(n.processed[msg.From], msg)
	n.vclock.TickUint64(msg.From)
	n.broadcastSnapshotState(context.Background(), msg.From, req.Token)
}

func (n *Node) broadcastSnapshotState(ctx context.Context, to uint64, token string) {
	defer n.snapshotMutex.Unlock()
	n.snapshotMutex.Lock()

	defer n.broadcastLock.Unlock()
	n.broadcastLock.Lock()

	vc := n.vclock.Copy()
	time, _ := vc.TimeUint64(n.networkInfo.ID)

	state := n.State()
	state.Token = token

	msg := &broadcast.Message{
		ID:     time,
		From:   n.networkInfo.ID,
		To:     to,
		VClock: vc,
		Type:   broadcast.TypeSnapshotRequest,
		Data:   state,
	}

	//defer n.processedMutex.Unlock()
	//n.processedMutex.Lock()

	n.processed[n.networkInfo.ID] = append(n.processed[n.networkInfo.ID], msg)
	n.broadcast.Broadcast(msg, n.connectedNodes)
}

func (n *Node) commitSnapshotState(msg *broadcast.Message) {
	st, ok := msg.Data.(*broadcast.State)
	if !ok {
		n.logger.Error(
			"couldn't commit snapshot stae, data type not snapshot state",
			"type", fmt.Sprintf("%T", msg.Data),
		)
	}

	n.processed[msg.From] = append(n.processed[msg.From], msg)
	n.vclock.TickUint64(msg.From)

	defer n.snapshotMutex.Unlock()
	n.snapshotMutex.Lock()

	n.snapshots[st.Token].AddState(st, msg.From)
	if n.snapshots[st.Token].Finished() {
		n.finishSnapshot(st.Token)
	}
}

func (n *Node) finishSnapshot(token string) {
	ss := n.snapshots[token]
	states, cErr := ss.GetStates()
	if cErr.IsNotNil() {
		n.logger.Fatal("finished snapshot not finished")
	}

	global := uint64(0)
	for _, v := range states {
		global += v.BitcakeBalance
		for ks, vs := range v.Sent {
			for _, smsg := range vs {
				has := false
				for _, rmsg := range states[ks].Recd[ks] {
					if rmsg.ID == smsg.ID {
						has = true
						break
					}
				}
				tx, ok := smsg.Data.(*broadcast.Transaction)
				if !ok {
					n.logger.Error("coudln't cast to transaction on finish snapshot")
					continue
				}
				if !has {
					n.logger.Info("[channel] bitcakes found", "bitcakes", tx.Bitcakes)
				}
			}
		}
	}

	n.logger.Info("[snapshot] global", "global", global)
}
