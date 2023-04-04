package app

import (
	"encoding/json"
	"fmt"
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
	app.Configuration = readTerraformExternalProtocol()
	app.debugEnabled = app.Configuration.Debug

	app.debug("Restoring application state")
	app.LoadState()

	app.debug("Updating runner assignments")
	var ci *gitlab.API = gitlab.NewAPI(app.GitLabHost, app.GitLabToken)
	var err error = ci.UpdateRunnerAssignments(app.RunnerTag)
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
	err = app.Save(app.ScalerState)
	if err != nil {
		log.Printf("Failed to save application state: %v", err)
	}
	writeTerraformExternalProtocol(&app.Fleet)
}

func (app *Application) LoadState() {
	var config = app.Configuration
	var err = app.LoadScalerState(config.ScalerState)
	if err != nil {
		log.Printf("Failed to load previous scaler state: %v", err)
	}
	err = app.LoadTerraformState(config.TerraformState)
	if err != nil {
		log.Printf("Failed to load terraform state: %v", err)
	}
	if config.RunnerAddress != "" {
		app.Entrypoint = config.RunnerAddress
	}
}

// Read configuration for Terraform external datasource (parse JSON from stdin)
func readTerraformExternalProtocol() Configuration {
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

// Write data for Terraform external data source (print JSON on stdout)
func writeTerraformExternalProtocol(fleet *cloud.Fleet) {
	var hosts = fleet.Hosts()
	var hostnames = make([]string, len(hosts))
	for i := 0; i < len(hosts); i++ {
		hostnames[i] = hosts[i].Name
	}
	var result = map[string]any{
		"runners": hostnames,
	}
	var output []byte
	var err error
	output, err = json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to compose Terraform output: %v", err)
	}
	fmt.Println(string(output))
}
