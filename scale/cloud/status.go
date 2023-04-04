package cloud

import (
	"fmt"
)

type HostStatus uint8

const New HostStatus = 0
const (
	Provisioning HostStatus = 1 << iota
	Ready
	Busy
	Idle
	OldAge
	Destroying
	Error
)

var StatusName = map[HostStatus]string{
	New:          "NEW",
	Provisioning: "PROVISIONING",
	Ready:        "READY",
	Busy:         "BUSY",
	Idle:         "IDLE",
	OldAge:       "OLDAGE",
	Destroying:   "DESTROYING",
	Error:        "ERROR",
}

func (status HostStatus) String() string {
	repr, ok := StatusName[status]
	if !ok {
		panic(fmt.Sprintf("invalid HostStatus=0b%b", status))
	}
	return repr
}

func (status HostStatus) Is(bitmask HostStatus) bool {
	if status == bitmask {
		return true
	}
	return status&bitmask != 0
}
