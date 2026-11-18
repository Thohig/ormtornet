package main

// En liten websocket-implementation (RFC 6455) med bara standardbiblioteket.
// Den stöder text- och binärmeddelanden, fragmentering, ping/pong och close.

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	wsGUID       = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	wsMaxMessage = 1 << 16
)

type wsConn struct {
	conn net.Conn
	br   *bufio.Reader
	wmu  sync.Mutex
}

func headerHas(h http.Header, name, token string) bool {
	for _, v := range h.Values(name) {
		for _, part := range strings.Split(v, ",") {
			if strings.EqualFold(strings.TrimSpace(part), token) {
				return true
			}
		}
	}
	return false
}

func wsUpgrade(w http.ResponseWriter, r *http.Request) (*wsConn, error) {
	if !headerHas(r.Header, "Connection", "upgrade") || !headerHas(r.Header, "Upgrade", "websocket") {
		http.Error(w, "websocket krävs", http.StatusBadRequest)
		return nil, errors.New("inte en websocket-förfrågan")
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		http.Error(w, "nyckel saknas", http.StatusBadRequest)
		return nil, errors.New("Sec-WebSocket-Key saknas")
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "kan inte uppgradera", http.StatusInternalServerError)
		return nil, errors.New("hijack stöds inte")
	}
	conn, brw, err := hj.Hijack()
	if err != nil {
		return nil, err
	}
	sum := sha1.Sum([]byte(key + wsGUID))
	accept := base64.StdEncoding.EncodeToString(sum[:])
	brw.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	brw.WriteString("Upgrade: websocket\r\nConnection: Upgrade\r\n")
	brw.WriteString("Sec-WebSocket-Accept: " + accept + "\r\n\r\n")
	if err := brw.Flush(); err != nil {
		conn.Close()
		return nil, err
	}
	return &wsConn{conn: conn, br: brw.Reader}, nil
}

// ReadMessage returnerar nästa hela data-meddelande.
func (c *wsConn) ReadMessage() ([]byte, error) {
	var msg []byte
	for {
		var hdr [2]byte
		if _, err := io.ReadFull(c.br, hdr[:]); err != nil {
			return nil, err
		}
		fin := hdr[0]&0x80 != 0
		op := hdr[0] & 0x0f
		masked := hdr[1]&0x80 != 0
		n := uint64(hdr[1] & 0x7f)
		switch n {
		case 126:
			var b [2]byte
			if _, err := io.ReadFull(c.br, b[:]); err != nil {
				return nil, err
			}
			n = uint64(binary.BigEndian.Uint16(b[:]))
		case 127:
			var b [8]byte
			if _, err := io.ReadFull(c.br, b[:]); err != nil {
				return nil, err
			}
			n = binary.BigEndian.Uint64(b[:])
		}
		if n > wsMaxMessage {
			return nil, errors.New("för stort meddelande")
		}
		var mask [4]byte
		if masked {
			if _, err := io.ReadFull(c.br, mask[:]); err != nil {
				return nil, err
			}
		}
		payload := make([]byte, n)
		if _, err := io.ReadFull(c.br, payload); err != nil {
			return nil, err
		}
		if masked {
			for i := range payload {
				payload[i] ^= mask[i%4]
			}
		}
		switch op {
		case 0x8: // close
			c.writeFrame(0x8, nil)
			return nil, io.EOF
		case 0x9: // ping
			c.writeFrame(0xA, payload)
		case 0xA: // pong
		case 0x0, 0x1, 0x2:
			msg = append(msg, payload...)
			if len(msg) > wsMaxMessage {
				return nil, errors.New("för stort meddelande")
			}
			if fin {
				return msg, nil
			}
		default:
			return nil, errors.New("okänd opcode")
		}
	}
}

func (c *wsConn) writeFrame(op byte, payload []byte) error {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	hdr := make([]byte, 0, 10)
	hdr = append(hdr, 0x80|op)
	n := len(payload)
	switch {
	case n < 126:
		hdr = append(hdr, byte(n))
	case n <= 0xffff:
		hdr = append(hdr, 126, byte(n>>8), byte(n))
	default:
		hdr = append(hdr, 127)
		hdr = binary.BigEndian.AppendUint64(hdr, uint64(n))
	}
	c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if _, err := c.conn.Write(append(hdr, payload...)); err != nil {
		return err
	}
	return nil
}

func (c *wsConn) WriteText(b []byte) error { return c.writeFrame(0x1, b) }

func (c *wsConn) Close() error { return c.conn.Close() }
