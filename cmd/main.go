package main

import (
	golog "log"

	"github.com/l2cup/kids2/internal/log"
	"github.com/l2cup/kids2/pkg/node"
)

func main() {
	logger, err := log.NewLogger(&log.Config{LogVerbosity: log.DebugVerbosity})
	if err != nil {
		golog.Fatal(err)
	}

	_, cErr := node.NewBootstrap(logger, "../config.json")
	if cErr.IsNotNil() {
		golog.Fatal(cErr)
	}

	//logger.Info("bs", "bootstrap", bs)
	select {}
}
