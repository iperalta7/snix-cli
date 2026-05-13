package console

import (
	"bytes"
	"errors"
	"net"
	"strings"
	"testing"
)

// TestWriteReadPacket_roundtrip verifies packet encoding survives a pipe.
func TestWriteReadPacket_roundtrip(t *testing.T) {
	cases := []rconPacket{
		{ID: 1, Type: rconAuth, Body: "mypassword"},
		{ID: 2, Type: rconExecCommand, Body: "say hello world"},
		{ID: 3, Type: rconResponse, Body: ""},
	}
	for _, want := range cases {
		client, server := net.Pipe()
		errCh := make(chan error, 1)
		var got rconPacket
		go func() {
			var err error
			got, err = readPacket(server)
			server.Close()
			errCh <- err
		}()
		if err := writePacket(client, want); err != nil {
			t.Fatalf("writePacket: %v", err)
		}
		client.Close()
		if err := <-errCh; err != nil {
			t.Fatalf("readPacket: %v", err)
		}
		if got.ID != want.ID || got.Type != want.Type || got.Body != want.Body {
			t.Errorf("roundtrip: got %+v, want %+v", got, want)
		}
	}
}

// fakeDialer returns a dialer that drives the server side in a goroutine.
func fakeDialer(responder func(conn net.Conn)) Dialer {
	return func(_ string) (net.Conn, error) {
		client, server := net.Pipe()
		go func() {
			defer server.Close()
			responder(server)
		}()
		return client, nil
	}
}

func TestRCONSend_success(t *testing.T) {
	var out bytes.Buffer
	r := &rconConsole{
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
			writePacket(conn, rconPacket{ID: 2, Type: rconResponse, Body: "players online: 3"}) //nolint:errcheck
		}),
	}
	if err := r.Send("list"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "players online: 3") {
		t.Errorf("output %q does not contain expected body", out.String())
	}
}

func TestRCONSend_authfail(t *testing.T) {
	r := &rconConsole{
		address:  "localhost:25575",
		password: "wrong",
		out:      &bytes.Buffer{},
		dial: fakeDialer(func(conn net.Conn) {
			readPacket(conn) //nolint:errcheck
			writePacket(conn, rconPacket{ID: rconAuthFail, Type: rconExecCommand}) //nolint:errcheck
		}),
	}
	err := r.Send("list")
	if err == nil {
		t.Fatal("expected auth failure error, got nil")
	}
	if !strings.Contains(err.Error(), "authentication failed") {
		t.Errorf("error %q does not mention authentication failure", err.Error())
	}
}

func TestRCONSend_dial_error(t *testing.T) {
	r := &rconConsole{
		address:  "localhost:25575",
		password: "secret",
		out:      &bytes.Buffer{},
		dial: func(_ string) (net.Conn, error) {
			return nil, errors.New("connection refused")
		},
	}
	err := r.Send("list")
	if err == nil {
		t.Fatal("expected dial error, got nil")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Errorf("error %q does not contain dial error", err.Error())
	}
	if !strings.Contains(err.Error(), "localhost:25575") {
		t.Errorf("error %q does not contain address", err.Error())
	}
}

func TestRCONAttach_error(t *testing.T) {
	r := &rconConsole{}
	err := r.Attach()
	if err == nil {
		t.Fatal("expected error from rcon Attach, got nil")
	}
	if !strings.Contains(err.Error(), "not supported") {
		t.Errorf("error %q does not mention not supported", err.Error())
	}
}
