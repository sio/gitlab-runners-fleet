package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

var DefaultConfiguration = Configuration{
	ScalerState:           "scaler.state",
	TerraformState:        "terraform.tfstate",
	RunnerTag:             "gitlab_runner_fleet",
	RunnerMaxJobs:         3,
	InstanceCountMax:      10,
	InstanceCountMin:      0,
	InstanceGrowMax:       3,
	InstanceProvisionTime: NewDuration(10 * 60 * second),
	InstanceMaxAge:        NewDuration(24 * 60 * 60 * second),
	InstanceMaxIdleTime:   NewDuration(40 * 60 * second),
}

const second = 1_000_000_000 // nanoseconds

type Configuration struct {
	GitLabHost            string   `json:"gitlab_host"`
	GitLabToken           string   `json:"gitlab_token"`
	ScalerState           string   `json:"scaler_state_file"`
	TerraformState        string   `json:"terraform_state_file"`
	RunnerAddress         string   `json:"runner_address"`
	RunnerTag             string   `json:"runner_tag"`
	RunnerMaxJobs         int      `json:"runner_max_jobs"`
	InstanceCountMax      int      `json:"instance_count_max"`
	InstanceCountMin      int      `json:"instance_count_min"`
	InstanceGrowMax       int      `json:"instance_grow_max"`
	InstanceProvisionTime Duration `json:"instance_provision_time"`
	InstanceMaxAge        Duration `json:"instance_max_age"`
	InstanceMaxIdleTime   Duration `json:"instance_max_idle_time"`
	Debug                 bool     `json:"debug"`
}

func (c *Configuration) ReadStdin() (err error) {
	err = json.NewDecoder(os.Stdin).Decode(c)
	return err
}

// Replace environment variable reference with its value
// if a corresponding JSON string starts with `env:` prefix
func (conf *Configuration) UnmarshalJSON(raw []byte) (err error) {
	var data = jsonConfigurationHolder(*conf)
	if err = json.Unmarshal(raw, &data); err != nil {
		return err
	}
	if data.GitLabHost, err = parseEnv(data.GitLabHost); err != nil {
		return fmt.Errorf("GitLabHost=%q: %w", data.GitLabHost, err)
	}
	if data.GitLabToken, err = parseEnv(data.GitLabToken); err != nil {
		return fmt.Errorf("GitLabToken=%q: %w", data.GitLabToken, err)
	}
	if data.ScalerState, err = parseEnv(data.ScalerState); err != nil {
		return fmt.Errorf("ScalerState=%q: %w", data.ScalerState, err)
	}
	if data.TerraformState, err = parseEnv(data.TerraformState); err != nil {
		return fmt.Errorf("TerraformState=%q: %w", data.TerraformState, err)
	}
	if data.RunnerTag, err = parseEnv(data.RunnerTag); err != nil {
		return fmt.Errorf("RunnerTag=%q: %w", data.RunnerTag, err)
	}
	*conf = Configuration(data)
	return nil
}

type jsonConfigurationHolder Configuration

func parseEnv(input string) (value string, err error) {
	const prefix = "env:"
	if !strings.HasPrefix(input, prefix) {
		return input, nil
	}
	var variable string = strings.TrimPrefix(input, prefix)
	value = os.Getenv(variable)
	if value == "" {
		return "", fmt.Errorf("variable not defined: %s", variable)
	}
	return value, nil
}

func NewDuration(value int64) (d Duration) {
	d.Duration = time.Duration(value)
	return d
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	var num int64
	num, err = json.Number(string(data)).Int64()
	var valueStr string
	if err == nil {
		valueStr = fmt.Sprintf("%ds", num) // unitless int is treated as number of seconds
	} else {
		valueStr = string(data[1 : len(data)-1]) // drop quote marks
	}
	d.Duration, err = time.ParseDuration(valueStr)
	return err
}
