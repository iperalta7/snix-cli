package config

import (
	"fmt"
	"regexp"
	"strings"
)

var validServiceName = regexp.MustCompile(`^[a-zA-Z0-9@._-]+$`)

// Validate checks all required fields and cross-field constraints.
// Collects all errors rather than failing on the first.
func (c *Config) Validate() error {
	var errs []string
	add := func(format string, args ...any) {
		errs = append(errs, fmt.Sprintf(format, args...))
	}

	// ── Top-level ─────────────────────────────────────────────────────────────
	if c.Name == "" {
		add("name is required")
	}
	if c.Runtime == "" {
		add("runtime is required")
	}

	switch c.Runtime {
	case "native":
		validateNative(c.Native, add)
	case "docker":
		validateDocker(c.Docker, add)
	case "":
		// already caught above
	default:
		add("runtime must be native or docker, got %q", c.Runtime)
	}

	// ── Service ───────────────────────────────────────────────────────────────
	if c.Service.Name == "" {
		add("service.name is required")
	} else if !validServiceName.MatchString(c.Service.Name) {
		add("service.name %q contains invalid characters (allowed: a-z A-Z 0-9 @ . _ -)", c.Service.Name)
	}

	// ── Console ───────────────────────────────────────────────────────────────
	switch c.Console.Type {
	case "rcon":
		if c.Console.RCON.Address == "" {
			add("console.rcon.address is required when console.type is rcon")
		}
		if c.Console.RCON.Password == "" {
			add("console.rcon.password or console.rcon.password_env is required when console.type is rcon")
		}
	case "session":
		if c.Runtime == "docker" {
			add("console.type session is not supported with docker runtime; use rcon")
		}
	case "":
		add("console.type is required (rcon or session)")
	default:
		add("console.type must be rcon or session, got %q", c.Console.Type)
	}

	// ── Backup ────────────────────────────────────────────────────────────────
	if c.Backup != nil {
		validateBackup(c.Backup, add)
	}

	// ── Update ────────────────────────────────────────────────────────────────
	if c.Update != nil {
		validateUpdate(c.Update, c.Runtime, add)
	}

	// ── Ports ─────────────────────────────────────────────────────────────────
	for i, p := range c.Ports {
		if p.Port < 1 || p.Port > 65535 {
			add("ports[%d]: port %d out of range (1-65535)", i, p.Port)
		}
		switch p.Protocol {
		case "tcp", "udp":
		case "":
			add("ports[%d]: protocol is required (tcp or udp)", i)
		default:
			add("ports[%d]: protocol must be tcp or udp, got %q", i, p.Protocol)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("config invalid:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func validateNative(n *NativeConfig, add func(string, ...any)) {
	if n == nil {
		add("native block is required when runtime is native")
		return
	}
	if n.Command == "" {
		add("native.command is required")
	}
	if n.Session == "" {
		add("native.session is required")
	}
	switch n.SessionType {
	case "screen", "tmux":
	case "":
		add("native.session_type is required (screen or tmux)")
	default:
		add("native.session_type must be screen or tmux, got %q", n.SessionType)
	}
	if n.StopTimeout <= 0 {
		add("native.stop_timeout must be greater than 0")
	}
}

func validateDocker(d *DockerConfig, add func(string, ...any)) {
	if d == nil {
		add("docker block is required when runtime is docker")
		return
	}
	if d.Image == "" {
		add("docker.image is required")
	}
	if d.ContainerName == "" {
		add("docker.container_name is required")
	}
	if d.StopTimeout <= 0 {
		add("docker.stop_timeout must be greater than 0")
	}
}

func validateBackup(b *BackupConfig, add func(string, ...any)) {
	if len(b.Dirs) == 0 && len(b.Files) == 0 {
		add("backup: at least one dir or file is required")
	}
	switch b.ArchiveFormat {
	case "tar.gz", "tar":
	case "":
		add("backup.archive_format is required (tar.gz or tar)")
	default:
		add("backup.archive_format must be tar.gz or tar, got %q", b.ArchiveFormat)
	}
	switch b.Backend {
	case "local":
		if b.Local == nil {
			add("backup.local block is required when backend is local")
		} else if b.Local.Path == "" {
			add("backup.local.path is required")
		}
	case "s3":
		if b.S3 == nil {
			add("backup.s3 block is required when backend is s3")
		} else {
			if b.S3.Bucket == "" {
				add("backup.s3.bucket is required")
			}
			if b.S3.Region == "" {
				add("backup.s3.region is required")
			}
		}
	case "oci":
		if b.OCI == nil {
			add("backup.oci block is required when backend is oci")
		} else {
			if b.OCI.BucketName == "" {
				add("backup.oci.bucket_name is required")
			}
			if b.OCI.Namespace == "" {
				add("backup.oci.namespace is required")
			}
			if b.OCI.Region == "" {
				add("backup.oci.region is required")
			}
		}
	case "gdrive":
		if b.GDrive == nil {
			add("backup.gdrive block is required when backend is gdrive")
		} else if b.GDrive.FolderID == "" {
			add("backup.gdrive.folder_id is required")
		}
	case "":
		add("backup.backend is required (local, s3, oci, or gdrive)")
	default:
		add("backup.backend must be local, s3, oci, or gdrive, got %q", b.Backend)
	}
}

func validateUpdate(u *UpdateConfig, runtime string, add func(string, ...any)) {
	switch u.Type {
	case "steamcmd":
		if runtime == "docker" {
			add("update.type steamcmd is not supported with docker runtime; use image")
		}
		if u.SteamCMD == nil {
			add("update.steamcmd block is required when update.type is steamcmd")
		} else {
			if u.SteamCMD.AppID == 0 {
				add("update.steamcmd.app_id is required")
			}
			if u.SteamCMD.InstallDir == "" {
				add("update.steamcmd.install_dir is required")
			}
		}
	case "script":
		if u.Script == nil || u.Script.Path == "" {
			add("update.script.path is required when update.type is script")
		}
	case "image":
		if runtime == "native" {
			add("update.type image is not supported with native runtime; use steamcmd or script")
		}
	case "none":
	case "":
		add("update.type is required (steamcmd, script, image, or none)")
	default:
		add("update.type must be steamcmd, script, image, or none, got %q", u.Type)
	}
}
