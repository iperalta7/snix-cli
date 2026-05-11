package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempYAML(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidYAML(t *testing.T) {
	yaml := `
name: test-server
game: minecraft
runtime: native
service:
  name: minecraft
console:
  type: session
native:
  command: /usr/bin/java
  session: mc
  session_type: screen
  stop_timeout: 30
`
	path := writeTempYAML(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "test-server" {
		t.Errorf("Name: got %q, want %q", cfg.Name, "test-server")
	}
	if cfg.Runtime != "native" {
		t.Errorf("Runtime: got %q, want %q", cfg.Runtime, "native")
	}
	if cfg.Service.Name != "minecraft" {
		t.Errorf("Service.Name: got %q, want %q", cfg.Service.Name, "minecraft")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	path := writeTempYAML(t, "name: [\nbad yaml")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}

func TestLoad_DockerRuntime(t *testing.T) {
	yaml := `
name: cs2
runtime: docker
service:
  name: cs2
console:
  type: rcon
  rcon:
    address: localhost:27015
    password: secret
docker:
  image: joedwards32/cs2
  container_name: cs2
  stop_timeout: 10
`
	path := writeTempYAML(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Docker == nil {
		t.Fatal("expected Docker config to be non-nil")
	}
	if cfg.Docker.Image != "joedwards32/cs2" {
		t.Errorf("Docker.Image: got %q, want %q", cfg.Docker.Image, "joedwards32/cs2")
	}
	if cfg.Console.RCON.Address != "localhost:27015" {
		t.Errorf("RCON.Address: got %q, want %q", cfg.Console.RCON.Address, "localhost:27015")
	}
}

func TestLoad_BackupAndUpdate(t *testing.T) {
	yaml := `
name: valheim
runtime: native
service:
  name: valheim
console:
  type: session
native:
  command: /usr/bin/valheim_server
  session: valheim
  session_type: screen
  stop_timeout: 20
backup:
  dirs:
    - saves
  archive_format: tar.gz
  backend: local
  local:
    path: /var/backups/valheim
    retention_days: 7
update:
  type: steamcmd
  steamcmd:
    app_id: 896660
    install_dir: /srv/valheim
    validate: true
`
	path := writeTempYAML(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Backup == nil {
		t.Fatal("expected Backup to be non-nil")
	}
	if cfg.Backup.Local.RetentionDays != 7 {
		t.Errorf("RetentionDays: got %d, want 7", cfg.Backup.Local.RetentionDays)
	}
	if cfg.Update == nil {
		t.Fatal("expected Update to be non-nil")
	}
	if cfg.Update.SteamCMD.AppID != 896660 {
		t.Errorf("AppID: got %d, want 896660", cfg.Update.SteamCMD.AppID)
	}
}
