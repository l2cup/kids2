package node

import (
	"fmt"
	nodepb "github.com/l2cup/kids2/internal/proto/node"
)

type NetworkInfo struct {
	ID     uint64 `json:"id"`
	IPAddr string `json:"ip_addr,omitempty"`
	Port   int32  `json:"port,omitempty"`
}

func (ni *NetworkInfo) Address() string {
	return fmt.Sprintf("%s:%d", ni.IPAddr, ni.Port)
}

func (ni *NetworkInfo) Proto() *nodepb.NetworkInfo {
	return &nodepb.NetworkInfo{
		Id:     ni.ID,
		IpAddr: ni.IPAddr,
		Port:   ni.Port,
	}
}

func NetworkInfoFromProto(ni *nodepb.NetworkInfo) *NetworkInfo {
	return &NetworkInfo{
		ID:     ni.GetId(),
		IPAddr: ni.GetIpAddr(),
		Port:   ni.GetPort(),
	}
}

func NetworkInfosFromProto(nispb ...*nodepb.NetworkInfo) []*NetworkInfo {
	nis := make([]*NetworkInfo, 0, len(nispb))
	for _, ni := range nispb {
		nis = append(nis, NetworkInfoFromProto(ni))
	}

	return nis
}

func NetworkInfosToProto(nis ...*NetworkInfo) []*nodepb.NetworkInfo {
	nispb := make([]*nodepb.NetworkInfo, 0, len(nis))
	for _, ni := range nis {
		nispb = append(nispb, ni.Proto())
	}

	return nispb
}
