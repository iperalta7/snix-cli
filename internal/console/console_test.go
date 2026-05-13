package console

import (
	"strings"
	"testing"

	"github.com/iperalta7/gsm/internal/config"
)

func sessionCfg() *config.Config {
	return &config.Config{
		Console: config.ConsoleConfig{Type: "session"},
		Native: &config.NativeConfig{
			Session:     "mc",
			SessionType: "screen",
		},
	}
}

func rconCfg() *config.Config {
	return &config.Config{
		Console: config.ConsoleConfig{
			Type: "rcon",
			RCON: config.RCONConfig{Address: "localhost:25575", Password: "secret"},
		},
	}
}

func TestNew_session(t *testing.T) {
	c := New(sessionCfg())
	if _, ok := c.(*sessionConsole); !ok {
		t.Fatalf("expected *sessionConsole, got %T", c)
	}
}

func TestNew_rcon(t *testing.T) {
	c := New(rconCfg())
	if _, ok := c.(*rconConsole); !ok {
		t.Fatalf("expected *rconConsole, got %T", c)
	}
}

func TestNew_rcon_attach_error(t *testing.T) {
	err := New(rconCfg()).Attach()
	if err == nil {
		t.Fatal("expected error from rcon Attach, got nil")
	}
	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("error %q does not contain %q", err.Error(), "not supported")
	}
}
