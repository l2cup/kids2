package broadcast

import (
	"context"
	"time"

	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/network"
)

type BroadcastFunc func(message *Message, node *network.Info)

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
	case TypeSnapshot:
		b.broadcast(message, nodes, b.broadcastSnapshot)
	}
}

func (b *Broadcast) broadcast(message *Message, nodes []*network.Info, fn BroadcastFunc) {
	for _, node := range nodes {
		go fn(message, node)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

func (b *Broadcast) broadcastSnapshot(message *Message, node *network.Info) {
	conn, cErr := network.DialGRPC(node)
	if cErr.IsNotNil() {
		b.logger.Error(
			"[broadcast] couldn't broadcast transaction",
			"err", cErr,
			"node", node,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := nodepb.NewNodeClient(conn)
	_, err := client.Snapshot(ctx, message.Proto())
	if err != nil {
		b.logger.Error(
			"[broadcast] transaction broadcast failed",
			"err", err,
		)
	}
}
