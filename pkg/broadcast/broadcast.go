package broadcast

import (
	"context"
	"math/rand"
	"time"

	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/network"
)

type broadcastFunc func(message *Message, node *network.Info)

type Broadcast struct {
	logger *log.Logger
}

func NewBroadcast(logger *log.Logger) *Broadcast {
	return &Broadcast{
		logger: logger,
	}
}

func (b *Broadcast) Broadcast(message *Message, nodes []*network.Info) {
	switch message.Type {
	case TypeTransaction:
		b.broadcast(message, nodes, b.broadcastTransaction)
	case TypeSnapshotState:
		b.broadcast(message, nodes, b.broadcastSnapshotState)
	case TypeSnapshotRequest:
		b.broadcast(message, nodes, b.broadcastSnapshotRequest)
	}
}

func (b *Broadcast) broadcast(message *Message, nodes []*network.Info, bf broadcastFunc) {
	for _, node := range nodes {
		go func(node *network.Info) {
			msRand := time.Duration(rand.Intn(3)*(rand.Intn(250)+rand.Intn(1200))) * time.Millisecond
			time.Sleep(msRand)
			bf(message, node)
		}(node)
	}
}

func (b *Broadcast) broadcastTransaction(message *Message, node *network.Info) {
	conn, cErr := network.DialGRPC(node)
	if cErr.IsNotNil() {
		b.logger.Error(
			"[broadcast] couldn't broadcast transaction",
			"err", cErr,
			"node", node,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	client := nodepb.NewNodeClient(conn)
	_, err := client.Transaction(ctx, message.Proto())
	if err != nil {
		b.logger.Error(
			"[broadcast] transaction broadcast failed",
			"err", err,
		)
	}
}

func (b *Broadcast) broadcastSnapshotState(message *Message, node *network.Info) {
	conn, cErr := network.DialGRPC(node)
	if cErr.IsNotNil() {
		b.logger.Error(
			"[broadcast] couldn't broadcast snapshot state",
			"err", cErr,
			"node", node,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	client := nodepb.NewNodeClient(conn)
	_, err := client.SendSnapshot(ctx, message.Proto())
	if err != nil {
		b.logger.Error(
			"[broadcast] snapshot broadcast failed",
			"err", err,
		)
	}
}

func (b *Broadcast) broadcastSnapshotRequest(message *Message, node *network.Info) {
	conn, cErr := network.DialGRPC(node)
	if cErr.IsNotNil() {
		b.logger.Error(
			"[broadcast] couldn't broadcast snapshot request",
			"err", cErr,
			"node", node,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	client := nodepb.NewNodeClient(conn)
	_, err := client.RequestSnapshot(ctx, message.Proto())
	if err != nil {
		b.logger.Error(
			"[broadcast] snapshot request broadcast failed",
			"err", err,
		)
	}
}
