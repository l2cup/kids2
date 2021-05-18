package node

import (
	"encoding/json"
	"os"

	"github.com/l2cup/kids2/internal/errors"
	"github.com/l2cup/kids2/pkg/network"
)

type topologyConfiguration struct {
	TotalNodeCount uint64 `json:"node_count"`
	Clique         bool   `json:"clique"`
	Nodes          []*struct {
		*network.Info
		ConnectedNodes []uint64 `json:"connected_nodes"`
	} `json:"nodes"`
}

func parseAndSetTopologyConfiguration(b *Bootstrap, path string) errors.Error {
	data, err := os.ReadFile(path)
	if err != nil {
		return errors.NewInternalError(
			"[bootstrap] couldn't open config file",
			errors.InternalServerError,
			err,
			"path", path,
		)
	}

	cfgFile := &topologyConfiguration{}
	err = json.Unmarshal(data, cfgFile)
	if err != nil {
		return errors.NewInternalError(
			"couldn't unmarshal json file to config",
			errors.JsonUnmarshalError,
			err,
		)
	}

	cfgs := make([]*network.Configuration, 0, len(cfgFile.Nodes))
	for i := range cfgs {
		cfgs[i].ConnectedNodes = make([]*network.Info, 0)
	}

	cfgMap := make(map[uint64]*network.Configuration, len(cfgs))
	for _, node := range cfgFile.Nodes {
		cfg := &network.Configuration{
			NetworkInfo: &network.Info{
				ID:     node.ID,
				IPAddr: node.IPAddr,
				Port:   node.Port,
			},
			TotalNodeCount: cfgFile.TotalNodeCount,
		}

		if _, ok := cfgMap[node.ID]; ok {
			return errors.New(
				"repeated id for node",
				errors.InternalServerError,
				"id", node.ID,
			)
		}
		cfgMap[node.ID] = cfg
		cfgs = append(cfgs, cfg)
	}

	if cfgFile.Clique == false {
		for i, node := range cfgFile.Nodes {
			for _, connected := range node.ConnectedNodes {
				if _, ok := cfgMap[connected]; !ok {
					return errors.New(
						"node doesn't exist",
						errors.InternalServerError,
						"id", connected,
					)
				}
				cfgs[i].ConnectedNodes = append(
					cfgs[i].ConnectedNodes,
					cfgMap[connected].NetworkInfo,
				)
			}
		}
	} else {
		for i := range cfgs {
			for _, jcfg := range cfgs {
				if jcfg.NetworkInfo.ID == cfgs[i].NetworkInfo.ID {
					continue
				}
				cfgs[i].ConnectedNodes = append(
					cfgs[i].ConnectedNodes,
					jcfg.NetworkInfo,
				)
			}
		}
	}

	b.NodeConfigurations = cfgs
	b.TotalNodeCount = cfgFile.TotalNodeCount
	return errors.Nil()
}
