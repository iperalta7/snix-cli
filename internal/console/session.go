package console

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/iperalta7/gsm/internal/config"
)

// SyscallExec is the shape of syscall.Exec — injectable for tests.
type SyscallExec func(argv0 string, argv []string, envv []string) error

type sessionConsole struct {
	session     string
	sessionType string
	exec        Executor
	sysExec     SyscallExec
	lookPath    func(string) (string, error)
}

func newSessionConsole(cfg *config.Config) *sessionConsole {
	return &sessionConsole{
		session:     cfg.Native.Session,
		sessionType: cfg.Native.SessionType,
		exec:        realExecutor,
		sysExec:     syscall.Exec,
		lookPath:    exec.LookPath,
	}
}

// Attach replaces the current process with the multiplexer attach command.
func (s *sessionConsole) Attach() error {
	var bin string
	var argv []string
	switch s.sessionType {
	case "screen":
		bin = "screen"
		argv = []string{"screen", "-r", s.session}
	case "tmux":
		bin = "tmux"
		argv = []string{"tmux", "attach", "-t", s.session}
	}
	path, err := s.lookPath(bin)
	if err != nil {
		return fmt.Errorf("finding %s binary: %w", bin, err)
	}
	if err := s.sysExec(path, argv, os.Environ()); err != nil {
		return fmt.Errorf("attaching to session %q: %w", s.session, err)
	}
	return nil
}

// Send delivers cmd to the running session without attaching.
func (s *sessionConsole) Send(cmd string) error {
	var name string
	var args []string
	switch s.sessionType {
	case "screen":
		name = "screen"
		args = []string{"-S", s.session, "-X", "stuff", cmd + "\n"}
	case "tmux":
		name = "tmux"
		args = []string{"send-keys", "-t", s.session, cmd, "Enter"}
	}
	out, err := s.exec(name, args...)
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg != "" {
			return fmt.Errorf("sending command to session %q: %w\n%s", s.session, err, msg)
		}
		return fmt.Errorf("sending command to session %q: %w", s.session, err)
	}
	return nil
}
