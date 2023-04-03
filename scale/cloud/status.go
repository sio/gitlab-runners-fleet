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
	Destroying
	Error
)

var StatusName = map[HostStatus]string{
	New:          "NEW",
	Provisioning: "PROVISIONING",
	Ready:        "READY",
	Busy:         "BUSY",
	Idle:         "IDLE",
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
