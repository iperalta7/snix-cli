package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name    string `yaml:"name"`
	Game    string `yaml:"game"`
	Runtime string `yaml:"runtime"` // native | docker

	Native  *NativeConfig `yaml:"native,omitempty"`
	Docker  *DockerConfig `yaml:"docker,omitempty"`
	Console ConsoleConfig `yaml:"console"`
	Service ServiceConfig `yaml:"service"`
	Backup  *BackupConfig `yaml:"backup,omitempty"`
	Update  *UpdateConfig `yaml:"update,omitempty"`
	Ports   []PortConfig  `yaml:"ports,omitempty"`
}

// ── Native runtime ────────────────────────────────────────────────────────────

type NativeConfig struct {
	Command     string   `yaml:"command"`
	Args        []string `yaml:"args,omitempty"`
	WorkingDir  string   `yaml:"working_dir"`
	User        string   `yaml:"user"`
	Session     string   `yaml:"session"`
	SessionType string   `yaml:"session_type"` // screen | tmux
	StopCommand string   `yaml:"stop_command"`
	StopTimeout int      `yaml:"stop_timeout"`
}

// ── Docker runtime ────────────────────────────────────────────────────────────

type DockerConfig struct {
	Image         string            `yaml:"image"`
	ContainerName string            `yaml:"container_name"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Env           map[string]string `yaml:"env,omitempty"`
	StopTimeout   int               `yaml:"stop_timeout"`
}

// ── Console ───────────────────────────────────────────────────────────────────

type ConsoleConfig struct {
	Type string     `yaml:"type"` // rcon | session
	RCON RCONConfig `yaml:"rcon"`
}

type RCONConfig struct {
	Address     string `yaml:"address"`
	Password    string `yaml:"password"`
	PasswordEnv string `yaml:"password_env"`
}

// ── Service ───────────────────────────────────────────────────────────────────

type ServiceConfig struct {
	Name string `yaml:"name"`
}

// ── Backup ────────────────────────────────────────────────────────────────────

type BackupConfig struct {
	FlushCommand  string         `yaml:"flush_command"`
	FlushWait     int            `yaml:"flush_wait"`
	Dirs          []string       `yaml:"dirs,omitempty"`
	Files         []string       `yaml:"files,omitempty"`
	ArchiveFormat string         `yaml:"archive_format"`
	Backend       string         `yaml:"backend"` // local | s3 | oci | gdrive
	Local         *LocalBackend  `yaml:"local,omitempty"`
	S3            *S3Backend     `yaml:"s3,omitempty"`
	OCI           *OCIBackend    `yaml:"oci,omitempty"`
	GDrive        *GDriveBackend `yaml:"gdrive,omitempty"`
}

type LocalBackend struct {
	Path          string `yaml:"path"`
	RetentionDays int    `yaml:"retention_days"`
}

type S3Backend struct {
	Bucket string `yaml:"bucket"`
	Prefix string `yaml:"prefix"`
	Region string `yaml:"region"`
}

type OCIBackend struct {
	BucketName   string `yaml:"bucket_name"`
	Namespace    string `yaml:"namespace"`
	Region       string `yaml:"region"`
	ObjectPrefix string `yaml:"object_prefix"`
	AuthMethod   string `yaml:"auth_method"`
}

type GDriveBackend struct {
	FolderID        string `yaml:"folder_id"`
	CredentialsFile string `yaml:"credentials_file"`
}

// ── Update ────────────────────────────────────────────────────────────────────

type UpdateConfig struct {
	Type     string          `yaml:"type"` // steamcmd | script | image | none
	SteamCMD *SteamCMDConfig `yaml:"steamcmd,omitempty"`
	Script   *ScriptConfig   `yaml:"script,omitempty"`
	Image    *ImageConfig    `yaml:"image,omitempty"`
}

type SteamCMDConfig struct {
	AppID      int    `yaml:"app_id"`
	InstallDir string `yaml:"install_dir"`
	Validate   bool   `yaml:"validate"`
}

type ScriptConfig struct {
	Path string `yaml:"path"`
}

type ImageConfig struct {
	Prune bool `yaml:"prune"`
}

// ── Ports ─────────────────────────────────────────────────────────────────────

type PortConfig struct {
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Protocol string `yaml:"protocol"` // tcp | udp
}

// ── Loader ────────────────────────────────────────────────────────────────────

// Load reads and parses the YAML config at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Console.RCON.PasswordEnv != "" {
		val := os.Getenv(cfg.Console.RCON.PasswordEnv)
		if val == "" {
			return nil, fmt.Errorf("env var %q (console.rcon.password_env) is not set or empty", cfg.Console.RCON.PasswordEnv)
		}
		cfg.Console.RCON.Password = val
	}
	return &cfg, nil
}
