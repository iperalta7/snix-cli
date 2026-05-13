package console

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/iperalta7/gsm/internal/config"
)

type dockerConsole struct {
	containerName string
	rcon          *rconConsole
	lookPath      func(string) (string, error)
	sysExec       SyscallExec
}

func newDockerConsole(cfg *config.Config) *dockerConsole {
	return &dockerConsole{
		containerName: cfg.Docker.ContainerName,
		rcon:          newRCONConsole(cfg),
		lookPath:      exec.LookPath,
		sysExec:       syscall.Exec,
	}
}

// Attach replaces the current process with docker attach, giving a live
// view of the container's PID-1 stdin/stdout (the game server console).
func (d *dockerConsole) Attach() error {
	path, err := d.lookPath("docker")
	if err != nil {
		return fmt.Errorf("finding docker binary: %w", err)
	}
	argv := []string{"docker", "attach", d.containerName}
	if err := d.sysExec(path, argv, os.Environ()); err != nil {
		return fmt.Errorf("attaching to container %q: %w", d.containerName, err)
	}
	return nil
}

// Send delivers cmd via RCON.
func (d *dockerConsole) Send(cmd string) error {
	return d.rcon.Send(cmd)
}
