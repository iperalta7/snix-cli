package console

import (
	"errors"
	"strings"
	"testing"
)

func fakeLookPath(bin string) (string, error) {
	return "/usr/bin/" + bin, nil
}

func TestSessionAttach_screen(t *testing.T) {
	var gotArgv0 string
	var gotArgv []string
	s := &sessionConsole{
		session:     "mc",
		sessionType: "screen",
		lookPath:    fakeLookPath,
		sysExec: func(argv0 string, argv []string, _ []string) error {
			gotArgv0 = argv0
			gotArgv = argv
			return nil
		},
	}
	if err := s.Attach(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(gotArgv0, "screen") {
		t.Errorf("argv0 %q does not end in 'screen'", gotArgv0)
	}
	if len(gotArgv) != 3 || gotArgv[1] != "-r" || gotArgv[2] != "mc" {
		t.Errorf("unexpected argv: %v", gotArgv)
	}
}

func TestSessionAttach_tmux(t *testing.T) {
	var gotArgv0 string
	var gotArgv []string
	s := &sessionConsole{
		session:     "mc",
		sessionType: "tmux",
		lookPath:    fakeLookPath,
		sysExec: func(argv0 string, argv []string, _ []string) error {
			gotArgv0 = argv0
			gotArgv = argv
			return nil
		},
	}
	if err := s.Attach(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(gotArgv0, "tmux") {
		t.Errorf("argv0 %q does not end in 'tmux'", gotArgv0)
	}
	if len(gotArgv) != 4 || gotArgv[1] != "attach" || gotArgv[2] != "-t" || gotArgv[3] != "mc" {
		t.Errorf("unexpected argv: %v", gotArgv)
	}
}

func TestSessionAttach_sysExec_error(t *testing.T) {
	s := &sessionConsole{
		session:     "mc",
		sessionType: "screen",
		lookPath:    fakeLookPath,
		sysExec: func(_ string, _ []string, _ []string) error {
			return errors.New("exec failed")
		},
	}
	err := s.Attach()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "exec failed") {
		t.Errorf("error %q does not wrap underlying error", err.Error())
	}
}

func TestSessionSend_screen(t *testing.T) {
	var gotName string
	var gotArgs []string
	s := &sessionConsole{
		session:     "mc",
		sessionType: "screen",
		exec: func(name string, args ...string) ([]byte, error) {
			gotName = name
			gotArgs = args
			return nil, nil
		},
	}
	if err := s.Send("say hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "screen" {
		t.Errorf("expected screen, got %q", gotName)
	}
	want := []string{"-S", "mc", "-X", "stuff", "say hello\n"}
	if len(gotArgs) != len(want) {
		t.Fatalf("args %v != %v", gotArgs, want)
	}
	for i := range want {
		if gotArgs[i] != want[i] {
			t.Errorf("arg[%d]: got %q, want %q", i, gotArgs[i], want[i])
		}
	}
}

func TestSessionSend_tmux(t *testing.T) {
	var gotName string
	var gotArgs []string
	s := &sessionConsole{
		session:     "mc",
		sessionType: "tmux",
		exec: func(name string, args ...string) ([]byte, error) {
			gotName = name
			gotArgs = args
			return nil, nil
		},
	}
	if err := s.Send("say hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "tmux" {
		t.Errorf("expected tmux, got %q", gotName)
	}
	want := []string{"send-keys", "-t", "mc", "say hello", "Enter"}
	if len(gotArgs) != len(want) {
		t.Fatalf("args %v != %v", gotArgs, want)
	}
	for i := range want {
		if gotArgs[i] != want[i] {
			t.Errorf("arg[%d]: got %q, want %q", i, gotArgs[i], want[i])
		}
	}
}

func TestSessionSend_error_with_output(t *testing.T) {
	s := &sessionConsole{
		session:     "mc",
		sessionType: "screen",
		exec: func(_ string, _ ...string) ([]byte, error) {
			return []byte("session not found"), errors.New("exit status 1")
		},
	}
	err := s.Send("say hello")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "session not found") {
		t.Errorf("error %q should contain executor output", err.Error())
	}
	if !strings.Contains(err.Error(), "exit status 1") {
		t.Errorf("error %q should wrap original error", err.Error())
	}
}
