package app

import (
	"time"

	"scale/cloud"
)

func (app *Application) UpdateStatus(host *cloud.Host) {
	var (
		err     error
		metrics cloud.Metrics
		now     = time.Now().UTC()
		zero    time.Time
	)
	host.UpdatedAt = now
	metrics, err = app.Metrics(host)
	switch {
	case err != nil:
		host.Status = cloud.Error
		if now.Sub(host.CreatedAt) < app.Configuration.InstanceProvisionTime.Duration {
			host.Status = cloud.Provisioning
		}
	case metrics.JobsRunning > 0:
		host.Status = cloud.Busy
		host.IdleSince = zero
	case metrics.JobsTotal > host.JobsDone:
		host.JobsDone = metrics.JobsTotal
		host.IdleSince = now
		host.Status = cloud.Ready
	case host.IdleSince != zero && now.Sub(host.IdleSince) > app.Configuration.InstanceMaxIdleTime.Duration:
		host.Status = cloud.Idle
	default:
		host.Status = cloud.Ready
		if host.IdleSince == zero {
			host.IdleSince = now
		}
	}
}
