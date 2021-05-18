package network

import (
	nodepb "github.com/l2cup/kids2/internal/proto/node"
)

type Configuration struct {
	NetworkInfo    *Info   `json:"network_info"`
	ConnectedNodes []*Info `json:"connected_nodes"`
	TotalNodeCount uint64  `json:"node_count"`
}

func (cfg *Configuration) Proto() *nodepb.Configuration {
	return &nodepb.Configuration{
		NetworkInfo:    cfg.NetworkInfo.Proto(),
		ConnectedNodes: InfosToProto(cfg.ConnectedNodes...),
		TotalNodeCount: cfg.TotalNodeCount,
	}
}

func ConfigurationFromProto(cfg *nodepb.Configuration) *Configuration {
	return &Configuration{
		NetworkInfo:    InfoFromProto(cfg.GetNetworkInfo()),
		ConnectedNodes: InfosFromProto(cfg.GetConnectedNodes()...),
		TotalNodeCount: cfg.GetTotalNodeCount(),
	}
}
