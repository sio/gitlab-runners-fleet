package app

import (
	"log"
	"sync"
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
		app.debug("Failed to fetch metrics for %s: %v", host, err)
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
				break
			}
		}
		if !found {
			break
		}
	}

	// Grow fleet if current capacity is not enough
	var countProvisioning int
	for _, host = range hosts {
		if host.Status == cloud.Provisioning {
			countProvisioning++
		}
	}
	var jobsPending = ci.CountPendingJobs()
	if jobsPending > 0 {
		app.debug("CI jobs currently pending: %d", jobsPending)
	}
	for jobsCapacity < jobsPending &&
		countProvisioning < app.InstanceGrowMax &&
		countHealthy < app.InstanceCountMax {
		host = app.AddHost()
		app.debug("Add host to the fleet: %s", host)
		jobsCapacity += app.RunnerMaxJobs
		countHealthy++
		countProvisioning++
	}

	// Remove idle instances
	var wg sync.WaitGroup
	app.debug("Triggering graceful cleanup for instances about to be deleted")
	for _, host = range hosts {
		if host.Status.Is(cloud.Idle | cloud.OldAge | cloud.Error) {
			wg.Add(1)
			go func(host *cloud.Host) {
				defer wg.Done()
				err := app.Cleanup(host)
				if err != nil {
					log.Printf("cleanup failed for %s: %v", host, err)
				} else {
					app.debug("..cleanup complete: %s", host)
				}
			}(host)
		}
	}
	wg.Wait()
	for _, host = range hosts {
		if host.Status.Is(cloud.Destroying | cloud.Error) {
			app.Delete(host)
		}
	}
}
