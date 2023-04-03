package app

import (
	"fmt"
	"log"
	"time"
)

func Run() {
	fmt.Println(tfExternalDatasourceConfig())
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
