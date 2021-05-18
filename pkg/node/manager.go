package node

import (
	"context"
	"math/rand"
	"time"

	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/broadcast"
	"github.com/l2cup/kids2/pkg/network"
	"github.com/l2cup/kids2/pkg/vc"
	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type Manager struct {
	logger        *log.Logger
	bootstrapInfo *network.Info
	broadcast     *broadcast.Broadcast
}

func NewManager(logger *log.Logger, bootstrapInfo *network.Info) *Manager {
	return &Manager{
		logger:        logger,
		bootstrapInfo: bootstrapInfo,
		broadcast:     broadcast.NewBroadcast(logger),
	}
}

func (m *Manager) NewNode() *Node {
	node := &Node{
		bitcakeBalance: uint64(1000 + rand.Intn(5000)),
		logger:         m.logger,
		bootstrapInfo:  m.bootstrapInfo,
		vclock:         vc.New(),
		sent:           make(map[uint64][]*broadcast.Message),
		recd:           make(map[uint64][]*broadcast.Message),
		broadcast:      m.broadcast,
	}

	m.mustRegisterNode(node)
	go m.startgRPCServer(node)

	m.logger.Debug("node debug", "node", node)
	return node
}

func (m *Manager) mustRegisterNode(node *Node) {
	conn, cErr := network.DialGRPC(m.bootstrapInfo)
	if cErr.IsNotNil() {
		m.logger.Fatal(
			"[fatal] couldn't dial grpc to register",
			"err", cErr,
			"bootstrap_info", m.bootstrapInfo,
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := nodepb.NewBootstrapClient(conn)
	response, err := client.Register(ctx, &emptypb.Empty{})
	if err != nil {
		m.logger.Fatal(
			"[fatal] couldn't register node",
			"err", err,
			"bootstrap_info", m.bootstrapInfo,
		)
	}

	m.configureNode(node, network.ConfigurationFromProto(response))
	m.logger.Debug(
		"[node] registered",
		"network_info", node.networkInfo,
		"node_count", node.totalNodeCount,
		"connected_nodes", node.connectedNodes,
	)
}

func (m *Manager) configureNode(node *Node, cfg *network.Configuration) {
	node.networkInfo = cfg.NetworkInfo
	node.totalNodeCount = cfg.TotalNodeCount
	node.connectedNodes = cfg.ConnectedNodes
}

func (m *Manager) startgRPCServer(node *Node) {
	srv := grpc.NewServer()
	nodepb.RegisterNodeServer(srv, node)

	m.logger.Info(
		"[node] started serving rpc on",
		"id", node.networkInfo.ID,
		"addr", node.networkInfo.Address(),
	)

	if cErr := network.ListenAndServeGRPC(node.networkInfo, srv); cErr.IsNotNil() {
		m.logger.Fatal("[fatal] boostrap failed to listen and serve grpc")
	}
}
