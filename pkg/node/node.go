package node

import (
	"sync"

	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"github.com/l2cup/kids2/pkg/network"
)

type Node struct {
	nodepb.UnimplementedNodeServer
	logger *log.Logger

	bitcakeBalance uint64
	mutex          sync.Mutex

	bootstrapInfo  *network.Info
	networkInfo    *network.Info
	connectedNodes []*network.Info
	totalNodeCount uint64
}
