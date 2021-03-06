package vc

import (
	vcpb "github.com/l2cup/kids2/internal/proto/vc"
	"strconv"
	"sync"
)

const (
	Backward int32 = -1
	Equal    int32 = 0
	Forward  int32 = 1
)

type VectorClock struct {
	mutex  sync.Mutex
	vclock map[string]uint64
}

func New() *VectorClock {
	return &VectorClock{
		mutex:  sync.Mutex{},
		vclock: make(map[string]uint64),
	}
}

func (vc *VectorClock) TimeUint64(ID uint64) (uint64, bool) {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()
	idString := strconv.FormatUint(ID, 10)

	time, ok := vc.vclock[idString]
	return time, ok
}

func (vc *VectorClock) Time(ID string) (uint64, bool) {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	time, ok := vc.vclock[ID]
	return time, ok
}

func (vc *VectorClock) Set(ID string, time uint64) {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	vc.vclock[ID] = time
}

func (vc *VectorClock) SetUint64(ID, time uint64) {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	idString := strconv.FormatUint(ID, 10)
	vc.vclock[idString] = time
}

func (vc *VectorClock) Tick(ID string) {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	vc.vclock[ID] += 1
}

func (vc *VectorClock) TickUint64(ID uint64) {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	idString := strconv.FormatUint(ID, 10)
	vc.vclock[idString] += 1
}

func (vc *VectorClock) Copy() *VectorClock {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	vcCopy := New()
	for k, v := range vc.vclock {
		vcCopy.vclock[k] = v
	}

	return vcCopy
}

func (vc *VectorClock) Map() map[string]uint64 {
	return vc.vclock
}

func (vc *VectorClock) IsDescendantOrEqual(other *VectorClock) bool {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	if len(other.vclock) > len(vc.vclock) {
		return false
	}

	is := true
	for k, v := range vc.vclock {
		otherVal := other.vclock[k]
		if otherVal > v {
			is = false
			break
		}
	}

	return is
}

func (vc *VectorClock) Proto() *vcpb.VectorClock {
	defer vc.mutex.Unlock()
	vc.mutex.Lock()

	entries := make([]*vcpb.Entry, 0, len(vc.vclock))
	for k, v := range vc.vclock {
		uint64Key, err := strconv.ParseUint(k, 10, 0)
		if err != nil {
			return nil
		}

		entry := &vcpb.Entry{
			Key:   uint64Key,
			Value: v,
		}
		entries = append(entries, entry)
	}

	return &vcpb.VectorClock{
		KvPairs: entries,
	}
}

func VectorClockFromProto(vclockpb *vcpb.VectorClock) *VectorClock {
	vc := New()
	for _, kvPair := range vclockpb.KvPairs {
		vc.SetUint64(kvPair.GetKey(), kvPair.GetValue())
	}

	return vc
}
