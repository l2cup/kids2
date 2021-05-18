package network

import (
	"fmt"

	networkpb "github.com/l2cup/kids2/internal/proto/network"
)

type Info struct {
	ID     uint64 `json:"id"`
	IPAddr string `json:"ip_addr,omitempty"`
	Port   int32  `json:"port,omitempty"`
}

func (ni *Info) Address() string {
	return fmt.Sprintf("%s:%d", ni.IPAddr, ni.Port)
}

func (ni *Info) Proto() *networkpb.Info {
	return &networkpb.Info{
		Id:     ni.ID,
		IpAddr: ni.IPAddr,
		Port:   ni.Port,
	}
}

func InfoFromProto(ni *networkpb.Info) *Info {
	return &Info{
		ID:     ni.GetId(),
		IPAddr: ni.GetIpAddr(),
		Port:   ni.GetPort(),
	}
}

func InfosFromProto(nispb ...*networkpb.Info) []*Info {
	nis := make([]*Info, 0, len(nispb))
	for _, ni := range nispb {
		nis = append(nis, InfoFromProto(ni))
	}

	return nis
}

func InfosToProto(nis ...*Info) []*networkpb.Info {
	nispb := make([]*networkpb.Info, 0, len(nis))
	for _, ni := range nis {
		nispb = append(nispb, ni.Proto())
	}

	return nispb
}
