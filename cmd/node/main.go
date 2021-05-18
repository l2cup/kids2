package main

import (
	"github.com/l2cup/kids2/internal/log"
	"github.com/l2cup/kids2/pkg/node"
	golog "log"
)

func main() {
	logger, err := log.NewLogger(&log.Config{LogVerbosity: log.DebugVerbosity})
	if err != nil {
		golog.Fatal(err)
	}

	nm := node.NewManager(logger, node.DefaultBootstrapInfo)

	for i := 0; i < 5; i++ {
		node := nm.NewNode()
		logger.Info("bitcakes", "bc", node.BitcakeBalance)
	}
}
