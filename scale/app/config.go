package app

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Configuration struct {
	GitLabHost            EnvString `json:"gitlab_host"`
	GitLabToken           EnvString `json:"gitlab_token"`
	RunnerTag             EnvString `json:"runner_tag"`
	RunnerMaxJobs         int       `json:"runner_max_jobs"`
	InstanceCountMax      int       `json:"instance_count_max"`
	InstanceCountMin      int       `json:"instance_count_min"`
	InstanceGrowMax       int       `json:"instance_grow_max"`
	InstanceProvisionTime Duration  `json:"instance_provision_time"`
	InstanceMaxAge        Duration  `json:"instance_max_age"`
}

// A string that will be unmarshalled from an environment variable
// if a corresponding JSON value starts with `env:` prefix
type EnvString string

func (str *EnvString) UnmarshalJSON(data []byte) (err error) {
	const prefix = "env:"
	var value string
	err = json.Unmarshal(data, &value)
	if err != nil {
		return err
	}
	*str = EnvString(value)
	if strings.HasPrefix(value, prefix) {
		var variable string = strings.TrimPrefix(value, prefix)
		*str = EnvString(os.Getenv(variable))
		if *str == "" {
			return fmt.Errorf("variable not defined: %s", variable)
		}
	}
	return nil
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
