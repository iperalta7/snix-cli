package console

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/iperalta7/gsm/internal/config"
)

func dockerCfg() *config.Config {
	return &config.Config{
		Runtime: "docker",
		Docker:  &config.DockerConfig{ContainerName: "mc"},
		Console: config.ConsoleConfig{
			Type: "rcon",
			RCON: config.RCONConfig{Address: "localhost:25575", Password: "secret"},
		},
	}
}

func TestNew_docker_returns_dockerConsole(t *testing.T) {
	c := New(dockerCfg())
	if _, ok := c.(*dockerConsole); !ok {
		t.Fatalf("expected *dockerConsole for docker runtime, got %T", c)
	}
}

func TestDockerAttach_success(t *testing.T) {
	var gotArgv0 string
	var gotArgv []string
	d := &dockerConsole{
		containerName: "mc",
		lookPath:      fakeLookPath,
		sysExec: func(argv0 string, argv []string, _ []string) error {
			gotArgv0 = argv0
			gotArgv = argv
			return nil
		},
	}
	if err := d.Attach(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(gotArgv0, "docker") {
		t.Errorf("argv0 %q does not end in 'docker'", gotArgv0)
	}
	if len(gotArgv) != 3 || gotArgv[1] != "attach" || gotArgv[2] != "mc" {
		t.Errorf("unexpected argv: %v", gotArgv)
	}
}

func TestDockerAttach_sysExec_error(t *testing.T) {
	d := &dockerConsole{
		containerName: "mc",
		lookPath:      fakeLookPath,
		sysExec: func(_ string, _ []string, _ []string) error {
			return errors.New("exec failed")
		},
	}
	err := d.Attach()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "exec failed") {
		t.Errorf("error %q does not wrap underlying error", err.Error())
	}
}

func TestDockerSend_delegates_to_rcon(t *testing.T) {
	var out bytes.Buffer
	d := &dockerConsole{
		containerName: "mc",
		rcon: &rconConsole{
			address:  "localhost:25575",
			password: "secret",
			out:      &out,
			dial: fakeDialer(func(conn net.Conn) {
				if _, err := readPacket(conn); err != nil {
					return
				}
				writePacket(conn, rconPacket{ID: 1, Type: rconExecCommand}) //nolint:errcheck
				if _, err := readPacket(conn); err != nil {
					return
				}
				writePacket(conn, rconPacket{ID: 2, Type: rconResponse, Body: "ok"}) //nolint:errcheck
			}),
		},
	}
	if err := d.Send("say hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "ok") {
		t.Errorf("output %q does not contain expected response", out.String())
	}
}
