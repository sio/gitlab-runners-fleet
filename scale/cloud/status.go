package cloud

type HostStatus uint8

const (
	New HostStatus = iota
	Provisioning
	Ready
	Busy
	Idle
	Destroying
	Error
)
