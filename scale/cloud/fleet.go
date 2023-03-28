package cloud

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Host struct {
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	IdleSince time.Time  `json:"idle_since"`
	Status    HostStatus `json:"status"`
}

type Fleet struct {
	Hosts      []Host `json:"hosts"`
	Entrypoint string `json:"entrypoint"`
}

func (fleet *Fleet) Load(filename string) error {
	var data []byte
	var err error
	data, err = os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, fleet)
}

func (fleet *Fleet) Save(filename string) error {
	var output []byte
	var err error
	output, err = json.MarshalIndent(*fleet, "", "  ")
	if err != nil {
		return err
	}

	var temp *os.File
	temp, err = os.CreateTemp(filepath.Dir(filename), filepath.Base(filename))
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())

	if _, err = temp.Write(output); err != nil {
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	if err = os.Rename(temp.Name(), filename); err != nil {
		return err
	}
	return nil
}
