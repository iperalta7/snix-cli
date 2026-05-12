package config

import (
	"strings"
	"testing"
)

func validNativeConfig() *Config {
	return &Config{
		Name:    "test-server",
		Runtime: "native",
		Native: &NativeConfig{
			Command:     "/usr/bin/java",
			Session:     "test",
			SessionType: "screen",
			StopTimeout: 30,
		},
		Service: ServiceConfig{Name: "test"},
		Console: ConsoleConfig{Type: "session"},
	}
}

func validDockerConfig() *Config {
	return &Config{
		Name:    "test-server",
		Runtime: "docker",
		Docker: &DockerConfig{
			Image:         "itzg/minecraft-server",
			ContainerName: "minecraft",
			StopTimeout:   30,
		},
		Service: ServiceConfig{Name: "test"},
		Console: ConsoleConfig{
			Type: "rcon",
			RCON: RCONConfig{Address: "localhost:25575", Password: "secret"},
		},
	}
}

func TestValidate_ValidNative(t *testing.T) {
	if err := validNativeConfig().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ValidDocker(t *testing.T) {
	if err := validDockerConfig().Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MissingName(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Name = ""
	assertError(t, cfg, "name is required")
}

func TestValidate_MissingRuntime(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Runtime = ""
	assertError(t, cfg, "runtime is required")
}

func TestValidate_InvalidRuntime(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Runtime = "podman"
	assertError(t, cfg, `runtime must be native or docker, got "podman"`)
}

func TestValidate_MissingServiceName(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Service.Name = ""
	assertError(t, cfg, "service.name is required")
}

func TestValidate_NativeBlockMissing(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Native = nil
	assertError(t, cfg, "native block is required")
}

func TestValidate_NativeCommandMissing(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Native.Command = ""
	assertError(t, cfg, "native.command is required")
}

func TestValidate_NativeSessionMissing(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Native.Session = ""
	assertError(t, cfg, "native.session is required")
}

func TestValidate_NativeSessionTypeInvalid(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Native.SessionType = "mux"
	assertError(t, cfg, `native.session_type must be screen or tmux, got "mux"`)
}

func TestValidate_NativeStopTimeoutZero(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Native.StopTimeout = 0
	assertError(t, cfg, "native.stop_timeout must be greater than 0")
}

func TestValidate_DockerBlockMissing(t *testing.T) {
	cfg := validDockerConfig()
	cfg.Docker = nil
	assertError(t, cfg, "docker block is required")
}

func TestValidate_DockerImageMissing(t *testing.T) {
	cfg := validDockerConfig()
	cfg.Docker.Image = ""
	assertError(t, cfg, "docker.image is required")
}

func TestValidate_DockerContainerNameMissing(t *testing.T) {
	cfg := validDockerConfig()
	cfg.Docker.ContainerName = ""
	assertError(t, cfg, "docker.container_name is required")
}

func TestValidate_ConsoleSessionWithDocker(t *testing.T) {
	cfg := validDockerConfig()
	cfg.Console = ConsoleConfig{Type: "session"}
	assertError(t, cfg, "console.type session is not supported with docker runtime")
}

func TestValidate_ConsoleRCONMissingAddress(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Console = ConsoleConfig{Type: "rcon", RCON: RCONConfig{}}
	assertError(t, cfg, "console.rcon.address is required")
}

func TestValidate_ConsoleMissingType(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Console = ConsoleConfig{}
	assertError(t, cfg, "console.type is required")
}

func TestValidate_PortOutOfRange(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Ports = []PortConfig{{Name: "game", Port: 0, Protocol: "tcp"}}
	assertError(t, cfg, "port 0 out of range")
}

func TestValidate_PortInvalidProtocol(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Ports = []PortConfig{{Name: "game", Port: 25565, Protocol: "udp6"}}
	assertError(t, cfg, `protocol must be tcp or udp, got "udp6"`)
}

func TestValidate_UpdateSteamCMDWithDocker(t *testing.T) {
	cfg := validDockerConfig()
	cfg.Update = &UpdateConfig{
		Type:     "steamcmd",
		SteamCMD: &SteamCMDConfig{AppID: 123, InstallDir: "/srv"},
	}
	assertError(t, cfg, "update.type steamcmd is not supported with docker runtime")
}

func TestValidate_UpdateImageWithNative(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Update = &UpdateConfig{Type: "image"}
	assertError(t, cfg, "update.type image is not supported with native runtime")
}

func TestValidate_UpdateSteamCMDMissingAppID(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Update = &UpdateConfig{
		Type:     "steamcmd",
		SteamCMD: &SteamCMDConfig{InstallDir: "/srv"},
	}
	assertError(t, cfg, "update.steamcmd.app_id is required")
}

func TestValidate_UpdateScriptMissingPath(t *testing.T) {
	cfg := validNativeConfig()
	cfg.Update = &UpdateConfig{Type: "script", Script: &ScriptConfig{}}
	assertError(t, cfg, "update.script.path is required")
}

func TestValidate_MultipleErrors(t *testing.T) {
	cfg := &Config{}
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected errors, got nil")
	}
	msg := err.Error()
	for _, want := range []string{"name is required", "runtime is required", "service.name is required"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error missing %q: %v", want, err)
		}
	}
}

func assertError(t *testing.T, cfg *Config, substr string) {
	t.Helper()
	err := cfg.Validate()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", substr)
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("error %q does not contain %q", err.Error(), substr)
	}
}
