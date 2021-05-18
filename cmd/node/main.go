package main

import (
	"context"
	golog "log"
	"time"

	"github.com/l2cup/kids2/internal/log"
	"github.com/l2cup/kids2/pkg/node"
)

func main() {
	logger, err := log.NewLogger(&log.Config{LogVerbosity: log.DebugVerbosity})
	if err != nil {
		golog.Fatal(err)
	}

	nm := node.NewManager(logger, node.DefaultBootstrapInfo)

	nodes := make([]*node.Node, 0, 5)
	for i := 0; i < 5; i++ {
		n := nm.NewNode()
		nodes = append(nodes, n)
	}

	time.Sleep(2 * time.Second)

	go func() {
		nodes[0].BroadcastTransaction(context.Background(), 4, 100)
		nodes[0].BroadcastTransaction(context.Background(), 4, 300)
		nodes[0].BroadcastTransaction(context.Background(), 4, 100)
	}()
	go func() {
		nodes[2].BroadcastTransaction(context.Background(), 3, 300)
		nodes[2].BroadcastTransaction(context.Background(), 3, 150)
		nodes[3].BroadcastTransaction(context.Background(), 2, 200)
		nodes[3].BroadcastTransaction(context.Background(), 1, 300)
	}()
	go nodes[1].BroadcastTransaction(context.Background(), 0, 100)
	go nodes[4].BroadcastTransaction(context.Background(), 3, 150)

	select {}
}
