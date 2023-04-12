package cloud

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/sio/coolname"
)

type Host struct {
	Name        string     `json:"name"`
	CreatedAt   time.Time  `json:"created_at"`
	IdleSince   time.Time  `json:"idle_since"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Status      HostStatus `json:"status"`
	JobsDone    int        `json:"jobs_done"`
	JobsRunning int        `json:"-"`
}

func (h *Host) String() string {
	return fmt.Sprintf("%s[%s]", h.Name, h.Status)
}

type Fleet struct {
	hosts      map[string]*Host
	Entrypoint string
}

// A static slice of all hosts at this point in time
func (fleet *Fleet) Hosts() []*Host {
	if len(fleet.hosts) > 100 {
		panic("This app was never meant to manage a fleet that large! Thorough code review is required, a full rewrite might be in order")
	}
	var hosts = make([]*Host, 0, len(fleet.hosts))
	for _, h := range fleet.hosts {
		hosts = append(hosts, h)
	}
	return hosts
}

// Create new host record
func (fleet *Fleet) AddHost() *Host {
	var name string
	for {
		var err error
		name, err = coolname.SlugN(2)
		if err != nil {
			continue
		}
		var exists bool
		_, exists = fleet.Get(name)
		if exists {
			continue
		}
		break
	}
	var host = &Host{
		Name:      name,
		Status:    New,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if fleet.hosts == nil {
		fleet.hosts = make(map[string]*Host)
	}
	fleet.hosts[host.Name] = host
	return host
}

// Delete host record
func (fleet *Fleet) Delete(host *Host) {
	delete(fleet.hosts, host.Name)
}

// Get host by name
func (fleet *Fleet) Get(name string) (host *Host, ok bool) {
	host, ok = fleet.hosts[name]
	return host, ok
}

// Load Hosts state from a file
func (fleet *Fleet) LoadScaleState(filename string) error {
	var data []byte
	var err error
	data, err = os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, fleet)
}

// Save Hosts state to a file
func (fleet *Fleet) Save(filename string) error {
	var output []byte
	var err error
	output, err = json.MarshalIndent(fleet, "", "  ")
	if err != nil {
		return err
	}

	var temp *os.File
	temp, err = os.CreateTemp(filepath.Dir(filename), filepath.Base(filename))
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(temp.Name()) }()

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

func (fleet *Fleet) MarshalJSON() ([]byte, error) {
	var serial = &serializableFleet{}
	serial.Pack(fleet)
	return json.Marshal(serial)
}

func (fleet *Fleet) UnmarshalJSON(data []byte) error {
	var serial = &serializableFleet{}
	var err error
	if err = json.Unmarshal(data, serial); err != nil {
		return err
	}
	serial.Unpack(fleet)
	return nil
}

// JSON-friendly representation on Fleet struct
type serializableFleet struct {
	Hosts      []*Host `json:"hosts"`
	Entrypoint string  `json:"entrypoint"`
}

func (s *serializableFleet) Pack(f *Fleet) {
	s.Entrypoint = f.Entrypoint
	s.Hosts = f.Hosts()
	sort.Slice(s.Hosts, func(i, j int) bool { return s.Hosts[i].Name < s.Hosts[j].Name })
}
func (s *serializableFleet) Unpack(f *Fleet) {
	f.Entrypoint = s.Entrypoint
	f.hosts = make(map[string]*Host)
	var h *Host
	for _, h = range s.Hosts {
		f.hosts[h.Name] = h
	}
}
