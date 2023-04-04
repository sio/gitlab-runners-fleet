package app

import (
	"log"
	"sync"
	"time"

	"scale/cloud"
	"scale/gitlab"
)

type Application struct {
	Configuration
	cloud.Fleet
	debugEnabled bool
}

func (app *Application) Run() {
	app.Configuration = tfExternalDatasourceConfig()
	app.debugEnabled = app.Configuration.Debug

	app.debug("Restoring application state")
	app.LoadState()

	app.debug("Updating runner assignments")
	var ci *gitlab.API = gitlab.NewAPI(string(app.GitLabHost), string(app.GitLabToken))
	var err error = ci.UpdateRunnerAssignments(string(app.RunnerTag))
	if err != nil {
		log.Printf("failed to update runner assignments: %v", err)
	}

	app.debug("Querying instance status")
	var wg sync.WaitGroup
	for _, host := range app.Hosts() {
		wg.Add(1)
		go func(host *cloud.Host) {
			defer wg.Done()
			app.UpdateStatus(host)
			app.debug("..status updated: %s", host)
		}(host)
	}
	wg.Wait()

	app.debug("Calculating scaling actions")
	app.Scale(ci)

	app.debug("Saving application state")
	err = app.Save(string(app.ScalerState))
	if err != nil {
		log.Printf("Failed to save application state: %v", err)
	}
}

func (app *Application) LoadState() {
	var config = app.Configuration
	var err = app.LoadScalerState(string(config.ScalerState))
	if err != nil {
		log.Printf("Failed to load previous scaler state: %v", err)
	}
	err = app.LoadTerraformState(string(config.TerraformState))
	if err != nil {
		log.Printf("Failed to load terraform state: %v", err)
	}
	if config.RunnerAddress != "" {
		app.Entrypoint = config.RunnerAddress
	}
}

// Read configuration for Terraform external datasource (parse JSON from stdin)
func tfExternalDatasourceConfig() Configuration {
	const timeout = 1 * time.Second

	var config Configuration = DefaultConfiguration
	var errors = make(chan error, 1)
	go func() {
		errors <- config.ReadStdin()
	}()
	select {
	case err := <-errors:
		if err != nil {
			log.Fatal("failed to load configuration: ", err)
		}
	case <-time.After(timeout):
		log.Fatal("configuration not received on stdin, exiting")
	}
	return config
}

func (app *Application) debug(msg string, values ...any) {
	if !app.debugEnabled {
		return
	}
	log.Printf(msg, values...)
}
