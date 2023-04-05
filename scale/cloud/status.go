package cloud

import (
	"fmt"
)

type HostStatus uint8

const New HostStatus = 0 // not deployed yet
const (
	// Deployed but not ready to accept jobs yet
	Provisioning HostStatus = 1 << iota

	// Ready to accept jobs but not currenty executing any
	Ready

	// Currently executing one or more jobs
	Busy

	// Has not been running any jobs for a while
	Idle

	// Has reached maximum allowed age
	OldAge

	// Cleanup completed, instance is about to be destroyed
	Destroying

	// Irrecoverable error, instance will be destroyed
	Error
)

var statusName = map[HostStatus]string{
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
	repr, ok := statusName[status]
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
