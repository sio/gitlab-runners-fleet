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
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	IdleSince time.Time  `json:"idle_since"`
	UpdatedAt time.Time  `json:"updated_at"`
	Status    HostStatus `json:"status"`
}

func (h *Host) String() string {
	return fmt.Sprintf("(%s@%d)", h.Name, h.Status)
}

type Fleet struct {
	Hosts      map[string]*Host
	Entrypoint string
}

// Create new host record
func (fleet *Fleet) New() *Host {
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
	fleet.Hosts[host.Name] = host
	return host
}

// Delete host record
func (fleet *Fleet) Delete(name string) {
}

// Get host by name
func (fleet *Fleet) Get(name string) (host *Host, ok bool) {
	host, ok = fleet.Hosts[name]
	return host, ok
}

// Load Hosts state from a file
func (fleet *Fleet) Load(filename string) error {
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
	s.Hosts = make([]*Host, 0, len(f.Hosts))
	var h *Host
	for _, h = range f.Hosts {
		s.Hosts = append(s.Hosts, h)
	}
	sort.Slice(s.Hosts, func(i, j int) bool { return s.Hosts[i].Name < s.Hosts[j].Name })
}
func (s *serializableFleet) Unpack(f *Fleet) {
	f.Entrypoint = s.Entrypoint
	f.Hosts = make(map[string]*Host)
	var h *Host
	for _, h = range s.Hosts {
		f.Hosts[h.Name] = h
	}
}
