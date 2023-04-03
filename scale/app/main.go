package app

import (
	"fmt"
	"log"
	"time"

	"scale/cloud"
)

type Application struct {
	Configuration
	cloud.Fleet
}

func (app *Application) Run() {
	app.Configuration = tfExternalDatasourceConfig()
	app.LoadState()
	for _, host := range app.Hosts() {
		fmt.Printf("Updating %s... ", host)
		app.UpdateStatus(host)
		fmt.Printf("Done: %s\n", host)
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
