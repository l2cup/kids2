package node

import (
	"context"
	"fmt"
	"net"

	"github.com/l2cup/kids2/internal/errors"
	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Configuration struct {
	NetworkInfo    *NetworkInfo   `json:"network_info"`
	ConnectedNodes []*NetworkInfo `json:"connected_nodes"`
	TotalNodeCount uint64         `json:"node_count"`
}

func (cfg *Configuration) Proto() *nodepb.Configuration {
	return &nodepb.Configuration{
		NetworkInfo:    cfg.NetworkInfo.Proto(),
		ConnectedNodes: NetworkInfosToProto(cfg.ConnectedNodes...),
		TotalNodeCount: cfg.TotalNodeCount,
	}
}

func ConfigurationFromProto(cfg *nodepb.Configuration) *Configuration {
	return &Configuration{
		NetworkInfo:    NetworkInfoFromProto(cfg.GetNetworkInfo()),
		ConnectedNodes: NetworkInfosFromProto(cfg.GetConnectedNodes()...),
		TotalNodeCount: cfg.GetTotalNodeCount(),
	}
}

type Bootstrap struct {
	nodepb.UnimplementedBootstrapServer

	NetworkInfo        *NetworkInfo
	TotalNodeCount     uint64
	bootstrapped       uint64
	NodeConfigurations []*Configuration
	logger             *log.Logger
}

func NewBootstrap(logger *log.Logger, cfgPath string) (*Bootstrap, errors.Error) {
	bootstrap := &Bootstrap{
		logger: logger,
		NetworkInfo: &NetworkInfo{
			ID:     0x8807,
			Port:   8007,
			IPAddr: "localhost",
		},
	}

	cErr := parseTopologyConfiguration(bootstrap, cfgPath)
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

	return nil, nil
}

func (b *Bootstrap) startgRPCServer() {
	listen, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d", b.NetworkInfo.IPAddr, b.NetworkInfo.ID),
	)
	if err != nil {
		b.logger.Fatal(
			"[bootstrap][fatal] couldn't listen on tcp",
			"err", err,
		)
	}

	s := grpc.NewServer()
	nodepb.RegisterBootstrapServer(s, b)

	b.logger.Info("[bootstrap] started serving rpc on", b.NetworkInfo.IPAddr, b.NetworkInfo.ID)
	if err := s.Serve(listen); err != nil {
		b.logger.Fatal("[bootstrap][fatal]failed to listen", "err", err)
	}
}
