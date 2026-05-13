package console

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/iperalta7/gsm/internal/config"
)

const (
	rconAuth        int32 = 3
	rconExecCommand int32 = 2
	rconResponse    int32 = 0
	rconAuthFail    int32 = -1
)

type rconPacket struct {
	ID   int32
	Type int32
	Body string
}

type rconConsole struct {
	address  string
	password string
	dial     Dialer
	out      io.Writer
}

func newRCONConsole(cfg *config.Config) *rconConsole {
	return &rconConsole{
		address:  cfg.Console.RCON.Address,
		password: cfg.Console.RCON.Password,
		dial:     realDialer,
		out:      os.Stdout,
	}
}

// Attach always errors — RCON has no interactive terminal to attach to.
func (r *rconConsole) Attach() error {
	return fmt.Errorf("console attach is not supported for rcon — use 'gsm cmd' to send a command")
}

// Send authenticates, sends cmd, and prints the server's response.
func (r *rconConsole) Send(cmd string) error {
	conn, err := r.dial(r.address)
	if err != nil {
		return fmt.Errorf("connecting to rcon %q: %w", r.address, err)
	}
	defer conn.Close()

	if err := writePacket(conn, rconPacket{ID: 1, Type: rconAuth, Body: r.password}); err != nil {
		return fmt.Errorf("sending rcon auth: %w", err)
	}
	resp, err := readPacket(conn)
	if err != nil {
		return fmt.Errorf("reading rcon auth response: %w", err)
	}
	if resp.ID == rconAuthFail {
		return fmt.Errorf("rcon authentication failed")
	}

	if err := writePacket(conn, rconPacket{ID: 2, Type: rconExecCommand, Body: cmd}); err != nil {
		return fmt.Errorf("sending rcon command: %w", err)
	}
	resp, err = readPacket(conn)
	if err != nil {
		return fmt.Errorf("reading rcon response: %w", err)
	}

	body := strings.TrimSpace(resp.Body)
	if body != "" {
		fmt.Fprintln(r.out, body)
	}
	return nil
}

// writePacket encodes a Valve RCON packet onto w.
//
// Wire format (all little-endian int32):
//
//	[ Size ][ ID ][ Type ][ Body bytes ][ 0x00 ][ 0x00 ]
//
// Size = len(Body) + 10 (4 ID + 4 Type + 1 body-null + 1 pad-null).
func writePacket(w io.Writer, pkt rconPacket) error {
	body := []byte(pkt.Body)
	size := int32(len(body) + 10)
	for _, v := range []any{size, pkt.ID, pkt.Type} {
		if err := binary.Write(w, binary.LittleEndian, v); err != nil {
			return err
		}
	}
	if _, err := w.Write(body); err != nil {
		return err
	}
	_, err := w.Write([]byte{0x00, 0x00})
	return err
}

// readPacket decodes one Valve RCON packet from r.
func readPacket(r io.Reader) (rconPacket, error) {
	var size int32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return rconPacket{}, fmt.Errorf("reading packet size: %w", err)
	}
	buf := make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		return rconPacket{}, fmt.Errorf("reading packet body: %w", err)
	}
	pkt := rconPacket{
		ID:   int32(binary.LittleEndian.Uint32(buf[0:4])),
		Type: int32(binary.LittleEndian.Uint32(buf[4:8])),
		Body: string(buf[8 : size-2]), // strip body-null + pad-null
	}
	return pkt, nil
}
