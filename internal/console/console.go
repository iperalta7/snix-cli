package console

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/iperalta7/gsm/internal/config"
)

// Console is the common interface for session and RCON backends.
type Console interface {
	// Attach replaces the current process with an interactive console session.
	// Returns an error for backends that do not support attachment.
	Attach() error

	// Send delivers cmd to the server without attaching.
	Send(cmd string) error
}

// Executor runs a named command and returns its combined output.
type Executor func(name string, args ...string) ([]byte, error)

// Dialer opens a TCP connection to addr.
type Dialer func(addr string) (net.Conn, error)

// New returns a Console backed by cfg.Console.Type.
func New(cfg *config.Config) Console {
	switch cfg.Console.Type {
	case "session":
		return newSessionConsole(cfg)
	case "rcon":
		if cfg.Runtime == "docker" {
			return newDockerConsole(cfg)
		}
		return newRCONConsole(cfg)
	default:
		return errConsole{t: cfg.Console.Type}
	}
}

func realExecutor(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}

func realDialer(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}

// errConsole is a fallback for unknown console types (unreachable after Validate).
type errConsole struct{ t string }

func (e errConsole) Attach() error       { return fmt.Errorf("unknown console type %q", e.t) }
func (e errConsole) Send(_ string) error { return fmt.Errorf("unknown console type %q", e.t) }
