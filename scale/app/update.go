package app

import (
	"log"
	"time"

	"scale/cloud"
	"scale/gitlab"
)

func (app *Application) UpdateStatus(host *cloud.Host) {
	var (
		err     error
		metrics cloud.Metrics
		now     = time.Now().UTC()
		zero    time.Time
	)
	host.UpdatedAt = now
	if host.CreatedAt == zero {
		host.CreatedAt = now
	}
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
		host.JobsRunning = metrics.JobsRunning
	case now.Sub(host.CreatedAt) > app.InstanceMaxAge.Duration:
		host.Status = cloud.OldAge
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

func (app *Application) Scale(ci *gitlab.API) {
	var (
		jobsCapacity, countHealthy int
		hosts                      = app.Hosts()
		host                       *cloud.Host
	)

	// Count healthy instances
	for _, host = range hosts {
		if host.Status.Is(cloud.Provisioning | cloud.Ready | cloud.Busy) {
			countHealthy++
			jobsCapacity += app.RunnerMaxJobs - host.JobsRunning
		}
	}

	// Keep some idle instances running if total count is below minimum
	for countHealthy < app.InstanceCountMin {
		var found bool
		for _, host = range hosts {
			if host.Status == cloud.Idle {
				found = true
				host.Status = cloud.Ready
				countHealthy++
				jobsCapacity += app.RunnerMaxJobs - host.JobsRunning
				break
			}
		}
		if !found {
			break
		}
	}

	// Mark surplus instances for removal
	for countHealthy > app.InstanceCountMax {
		var found bool
		for _, host = range hosts {
			if host.Status.Is(cloud.Provisioning | cloud.Ready) {
				found = true
				host.Status = cloud.Idle
				countHealthy--
				jobsCapacity -= app.RunnerMaxJobs - host.JobsRunning
			}
		}
		if !found {
			break
		}
	}

	// Grow fleet if current capacity is not enough
	var jobsPending = ci.CountPendingJobs()
	for jobsCapacity < jobsPending {
		_ = app.AddHost()
		jobsCapacity += app.RunnerMaxJobs
		countHealthy++
	}

	// Remove idle instances
	for _, host = range hosts {
		if host.Status.Is(cloud.Idle | cloud.OldAge) {
			err := app.Cleanup(host)
			if err != nil {
				log.Printf("cleanup failed for %s: %v", host, err)
			}
		}
		if host.Status.Is(cloud.Destroying | cloud.Error) {
			app.Delete(host)
		}
	}
}
