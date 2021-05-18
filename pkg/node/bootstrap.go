package node

import (
	"context"

	"github.com/l2cup/kids2/internal/errors"
	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/network"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var DefaultBootstrapInfo = &network.Info{
	ID:     0x8807,
	Port:   8007,
	IPAddr: "localhost",
}

type Bootstrap struct {
	nodepb.UnimplementedBootstrapServer

	NetworkInfo        *network.Info
	TotalNodeCount     uint64
	bootstrapped       uint64
	NodeConfigurations []*network.Configuration
	logger             *log.Logger
}

func NewBootstrap(logger *log.Logger, cfgPath string) (*Bootstrap, errors.Error) {
	bootstrap := &Bootstrap{
		logger:      logger,
		NetworkInfo: DefaultBootstrapInfo,
	}

	cErr := parseAndSetTopologyConfiguration(bootstrap, cfgPath)
	if cErr.IsNotNil() {
		return nil, cErr
	}

	go bootstrap.startgRPCServer()

	return bootstrap, errors.Nil()
}

func (b *Bootstrap) Register(
	ctx context.Context,
	_ *emptypb.Empty,
) (*nodepb.Configuration, error) {
	if b.bootstrapped == b.TotalNodeCount {
		return nil, pkgerrors.New("all nodes bootstrapped")
	}

	cfg := b.NodeConfigurations[b.bootstrapped].Proto()
	b.bootstrapped++
	b.logger.Debug(
		"[bootstrap] node bootstrapped",
		"node_configuration", cfg,
	)
	return cfg, nil
}

func (b *Bootstrap) startgRPCServer() {
	srv := grpc.NewServer()
	nodepb.RegisterBootstrapServer(srv, b)

	b.logger.Info(
		"[bootstrap] started serving rpc on",
		"addr", b.NetworkInfo.Address(),
	)

	if cErr := network.ListenAndServeGRPC(b.NetworkInfo, srv); cErr.IsNotNil() {
		b.logger.Fatal("[fatal] boostrap failed to listen and serve grpc")
	}
}
