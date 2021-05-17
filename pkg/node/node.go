package node

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/l2cup/kids2/internal/log"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
	"google.golang.org/grpc"
)

type Bitcake = uint64

type Node struct {
	NetworkInfo    *NetworkInfo `json:"network_info"`
	ConnectedNodes []*Node      `json:"connected_nodes"`
	Bitcakes       Bitcake      `json:"bitcakes"`
	Valid          bool         `json:"valid"`
	bootstrapInfo  *NetworkInfo
}

func NewNode(logger *log.Logger, bootstrapInfo *NetworkInfo) *Node {
	node := &Node{
		Bitcakes: uint64(1000 + rand.Intn(5000)),
	}

	addr := fmt.Sprintf("%s:%d", bootstrapInfo.IPAddr, bootstrapInfo.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.Fatal("err", "err", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := nodepb.NewBootstrapClient(conn)
	r, err := client.Register(ctx, nil)
	if err != nil {
		logger.Fatal("couldn't register fuck", "err", err)
	}

	logger.Info("registered", r)

	return node
}
