package process

import (
	"errors"
	"testing"
)

func mockExec(outputs map[string][]byte, errs map[string]error) Executor {
	return func(name string, args ...string) ([]byte, error) {
		key := name
		if len(args) > 0 {
			key = args[0]
		}
		out := outputs[key]
		err := errs[key]
		return out, err
	}
}

func TestStart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := &Manager{exec: mockExec(map[string][]byte{"start": nil}, nil)}
		if err := m.Start("myserver"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("systemctl error with output", func(t *testing.T) {
		m := &Manager{exec: mockExec(
			map[string][]byte{"start": []byte("Unit myserver.service not found.")},
			map[string]error{"start": errors.New("exit status 5")},
		)}
		err := m.Start("myserver")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestStop(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := &Manager{exec: mockExec(map[string][]byte{"stop": nil}, nil)}
		if err := m.Stop("myserver"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error propagated", func(t *testing.T) {
		m := &Manager{exec: mockExec(nil, map[string]error{"stop": errors.New("exit status 1")})}
		if err := m.Stop("myserver"); err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestRestart(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := &Manager{exec: mockExec(map[string][]byte{"restart": nil}, nil)}
		if err := m.Restart("myserver"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGetStatus(t *testing.T) {
	cases := []struct {
		name        string
		output      string
		wantActive  bool
		wantState   string
		wantSubState string
	}{
		{
			name:         "active running",
			output:       "ActiveState=active\nSubState=running\n",
			wantActive:   true,
			wantState:    "active",
			wantSubState: "running",
		},
		{
			name:         "inactive dead",
			output:       "ActiveState=inactive\nSubState=dead\n",
			wantActive:   false,
			wantState:    "inactive",
			wantSubState: "dead",
		},
		{
			name:         "failed",
			output:       "ActiveState=failed\nSubState=failed\n",
			wantActive:   false,
			wantState:    "failed",
			wantSubState: "failed",
		},
		{
			name:         "activating",
			output:       "ActiveState=activating\nSubState=start\n",
			wantActive:   false,
			wantState:    "activating",
			wantSubState: "start",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := &Manager{exec: func(name string, args ...string) ([]byte, error) {
				return []byte(tc.output), nil
			}}
			s, err := m.GetStatus("myserver")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if s.Active != tc.wantActive {
				t.Errorf("Active: got %v, want %v", s.Active, tc.wantActive)
			}
			if s.State != tc.wantState {
				t.Errorf("State: got %q, want %q", s.State, tc.wantState)
			}
			if s.SubState != tc.wantSubState {
				t.Errorf("SubState: got %q, want %q", s.SubState, tc.wantSubState)
			}
		})
	}

	t.Run("empty output returns error", func(t *testing.T) {
		m := &Manager{exec: func(name string, args ...string) ([]byte, error) {
			return []byte("SubState=dead\n"), nil
		}}
		_, err := m.GetStatus("myserver")
		if err == nil {
			t.Fatal("expected error for missing ActiveState, got nil")
		}
	})

	t.Run("executor error", func(t *testing.T) {
		m := &Manager{exec: func(name string, args ...string) ([]byte, error) {
			return nil, errors.New("exit status 4")
		}}
		_, err := m.GetStatus("myserver")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestParseProperties(t *testing.T) {
	out := "ActiveState=active\nSubState=running\n"
	props := parseProperties(out)
	if props["ActiveState"] != "active" {
		t.Errorf("ActiveState: got %q, want %q", props["ActiveState"], "active")
	}
	if props["SubState"] != "running" {
		t.Errorf("SubState: got %q, want %q", props["SubState"], "running")
	}
}
