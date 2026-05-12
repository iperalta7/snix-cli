package process

import (
	"fmt"
	"os/exec"
	"strings"
)

// Executor runs a named command and returns its combined output.
type Executor func(name string, args ...string) ([]byte, error)

// Manager controls a game server process via systemd.
type Manager struct {
	exec Executor
}

// New returns a Manager backed by the real systemctl binary.
func New() *Manager {
	return &Manager{exec: realExecutor}
}

func realExecutor(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

// Status holds the parsed state of a systemd unit.
type Status struct {
	Service  string
	Active   bool
	State    string // active | inactive | failed | activating | deactivating
	SubState string // running | dead | exited | ...
}

func (s *Status) String() string {
	return fmt.Sprintf("service:  %s\nstate:    %s (%s)", s.Service, s.State, s.SubState)
}

// Start runs `systemctl start <service>`.
func (m *Manager) Start(service string) error {
	if err := m.systemctl("start", service); err != nil {
		return fmt.Errorf("starting %s: %w", service, err)
	}
	return nil
}

// Stop runs `systemctl stop <service>`.
func (m *Manager) Stop(service string) error {
	if err := m.systemctl("stop", service); err != nil {
		return fmt.Errorf("stopping %s: %w", service, err)
	}
	return nil
}

// Restart runs `systemctl restart <service>`.
func (m *Manager) Restart(service string) error {
	if err := m.systemctl("restart", service); err != nil {
		return fmt.Errorf("restarting %s: %w", service, err)
	}
	return nil
}

// GetStatus returns the current state of the systemd unit.
func (m *Manager) GetStatus(service string) (*Status, error) {
	out, err := m.exec("systemctl", "show", "--property=ActiveState,SubState", service)
	if err != nil {
		return nil, fmt.Errorf("querying status of %s: %w", service, err)
	}

	props := parseProperties(string(out))
	state := props["ActiveState"]
	sub := props["SubState"]

	if state == "" {
		return nil, fmt.Errorf("could not determine state for service %q", service)
	}

	return &Status{
		Service:  service,
		Active:   state == "active",
		State:    state,
		SubState: sub,
	}, nil
}

func (m *Manager) systemctl(verb, service string) error {
	out, err := m.exec("systemctl", verb, service)
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("%w\n%s", err, msg)
		}
		return err
	}
	return nil
}

func parseProperties(output string) map[string]string {
	props := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		k, v, ok := strings.Cut(line, "=")
		if ok {
			props[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
	}
	return props
}
