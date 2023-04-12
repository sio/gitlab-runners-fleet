package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// Load current Terraform state
func (fleet *Fleet) LoadTerraformState(statefile string) (err error) {
	err = fleet.loadTerraformStateCLI()
	if err == nil {
		return nil
	}
	return fleet.loadTerraformStateFile(statefile)
}

// Load Terraform state by querying CLI (`terraform show`)
func (fleet *Fleet) loadTerraformStateCLI() (err error) {
	var (
		tfstate any
		tfjson  []byte
	)
	tfjson, err = exec.Command("terraform", "show", "--json").Output()
	if err != nil {
		return err
	}
	err = json.Unmarshal(tfjson, &tfstate)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	var entrypoint string
	entrypoint, err = jsGet[string](tfstate, "values", "outputs", "external_ip", "value")
	if err != nil || entrypoint == "" {
		return fmt.Errorf("entrypoint not found in terraform output")
	}
	fleet.Entrypoint = entrypoint

	resources, err := jsGet[[]any](tfstate, "values", "root_module", "resources")
	if err != nil {
		return fmt.Errorf("failed to read terraform resources")
	}
	if fleet.hosts == nil {
		fleet.hosts = make(map[string]*Host)
	}
	for _, r := range resources {
		var resourceType string
		resourceType, err = jsGet[string](r, "type")
		if err != nil || resourceType != "yandex_compute_instance" {
			continue
		}
		var host = &Host{}
		host.Name, err = jsGet[string](r, "values", "name")
		if err != nil || host.Name == "" || host.Name == "gateway" {
			continue
		}
		var ctime string
		ctime, err = jsGet[string](r, "values", "created_at")
		if err == nil {
			t, err := time.Parse(time.RFC3339, ctime)
			if err == nil {
				host.CreatedAt = t
			}
		}
		_, exists := fleet.hosts[host.Name]
		if exists {
			continue
		}
		fleet.hosts[host.Name] = host
	}
	return nil
}

// Load Terraform state by reading state file directly
func (fleet *Fleet) loadTerraformStateFile(filename string) (err error) {
	var (
		tfstate any
		tfjson  []byte
	)
	tfjson, err = os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tfjson, &tfstate)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	var entrypoint string
	entrypoint, err = jsGet[string](tfstate, "outputs", "external_ip", "value")
	if err != nil || entrypoint == "" {
		return fmt.Errorf("terraform state file does not contain previous value for external_ip")
	}
	fleet.Entrypoint = entrypoint

	resources, err := jsGet[[]any](tfstate, "resources")
	if err != nil {
		return fmt.Errorf("failed to read terraform resources")
	}
	if fleet.hosts == nil {
		fleet.hosts = make(map[string]*Host)
	}
	for _, r := range resources {
		var resourceType string
		resourceType, err = jsGet[string](r, "type")
		if err != nil || resourceType != "yandex_compute_instance" {
			continue
		}
		var instances []any
		instances, err = jsGet[[]any](r, "instances")
		if err != nil {
			continue
		}
		for _, i := range instances {
			var host = &Host{}
			host.Name, err = jsGet[string](i, "attributes", "name")
			if err != nil || host.Name == "" || host.Name == "gateway" {
				continue
			}
			var ctime string
			ctime, err = jsGet[string](i, "attributes", "created_at")
			if err == nil {
				t, err := time.Parse(time.RFC3339, ctime)
				if err == nil {
					host.CreatedAt = t
				}
			}
			_, exists := fleet.hosts[host.Name]
			if exists {
				continue
			}
			fleet.hosts[host.Name] = host
		}
	}
	return nil
}
