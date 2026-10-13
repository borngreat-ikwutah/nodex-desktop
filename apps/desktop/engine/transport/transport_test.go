package transport

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"nodex-desktop/engine/identity"
)

// chunkedReader simulates fragmented TCP packet delivery by yielding 1 byte at a time.
type chunkedReader struct {
	r io.Reader
}

func (c *chunkedReader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	return c.r.Read(p[:1])
}

func TestProtocolFramingRoundTrip(t *testing.T) {
	orig := &Message{
		Type:      TypePing,
		SenderID:  "111 111 111",
		TargetID:  "222 222 222",
		UUID:      "uuid-1234",
		PublicKey: "pubkey-abcd",
		Timestamp: 123456789,
		Payload:   json.RawMessage(`{"hello":"world"}`),
	}

	var buf bytes.Buffer
	if err := WriteMessage(&buf, orig); err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}

	// Verify header length prefix matches payload length
	if buf.Len() < HeaderSize {
		t.Fatalf("buffer shorter than header size")
	}
	payloadLen := binary.BigEndian.Uint32(buf.Bytes()[:HeaderSize])
	if int(payloadLen) != buf.Len()-HeaderSize {
		t.Fatalf("header length %d != actual payload length %d", payloadLen, buf.Len()-HeaderSize)
	}

	// Read message back
	parsed, err := ReadMessage(&buf)
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}

	if parsed.Type != orig.Type || parsed.SenderID != orig.SenderID || parsed.UUID != orig.UUID {
		t.Fatalf("mismatch: got %+v, want %+v", parsed, orig)
	}
	if string(parsed.Payload) != string(orig.Payload) {
		t.Fatalf("payload mismatch: got %s, want %s", string(parsed.Payload), string(orig.Payload))
	}
}

func TestProtocolFramingFragmentedStream(t *testing.T) {
	orig := &Message{
		Type:     TypePing,
		SenderID: "333 333 333",
		Payload:  json.RawMessage(`{"stream":"fragmented"}`),
	}

	var buf bytes.Buffer
	if err := WriteMessage(&buf, orig); err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}

	// Wrap in chunked reader delivering 1 byte at a time
	cr := &chunkedReader{r: &buf}
	parsed, err := ReadMessage(cr)
	if err != nil {
		t.Fatalf("ReadMessage over fragmented stream failed: %v", err)
	}

	if parsed.SenderID != orig.SenderID {
		t.Fatalf("expected SenderID %q, got %q", orig.SenderID, parsed.SenderID)
	}
}

func TestProtocolOversizedMessageRejected(t *testing.T) {
	var buf bytes.Buffer
	var header [HeaderSize]byte
	// 17 MB length prefix (exceeds MaxPayloadSize of 16MB)
	binary.BigEndian.PutUint32(header[:], 17*1024*1024)
	buf.Write(header[:])

	_, err := ReadMessage(&buf)
	if !errors.Is(err, ErrMessageTooLarge) {
		t.Fatalf("expected ErrMessageTooLarge, got %v", err)
	}
}

func TestServerClientPingPong(t *testing.T) {
	serverID := identity.Public{
		NodeID:    "222 222 222",
		UUID:      "server-uuid",
		PublicKey: "server-pubkey",
	}
	clientID := identity.Public{
		NodeID:    "111 111 111",
		UUID:      "client-uuid",
		PublicKey: "client-pubkey",
	}

	// Start server on ephemeral port (port 0)
	srv := NewServer(serverID,
		WithServerPort(0),
		WithServerListenIP("127.0.0.1"),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := srv.Start(ctx); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer srv.Stop()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", srv.Port())

	// Client pings server
	client := NewClient(clientID, 2*time.Second)
	result, err := client.Ping(ctx, serverAddr)
	if err != nil {
		t.Fatalf("client.Ping failed: %v", err)
	}

	if result.PeerID != serverID.NodeID {
		t.Errorf("expected PeerID %q, got %q", serverID.NodeID, result.PeerID)
	}
	if result.PeerUUID != serverID.UUID {
		t.Errorf("expected PeerUUID %q, got %q", serverID.UUID, result.PeerUUID)
	}
	if result.PeerPublicKey != serverID.PublicKey {
		t.Errorf("expected PeerPublicKey %q, got %q", serverID.PublicKey, result.PeerPublicKey)
	}
	if result.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", result.Latency)
	}
}

func TestConcurrentClientPings(t *testing.T) {
	serverID := identity.Public{
		NodeID: "999 999 999",
		UUID:   "srv-concurrent-uuid",
	}
	srv := NewServer(serverID, WithServerPort(0), WithServerListenIP("127.0.0.1"))
	if err := srv.Start(context.Background()); err != nil {
		t.Fatalf("server start: %v", err)
	}
	defer srv.Stop()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", srv.Port())

	var wg sync.WaitGroup
	errCh := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cid := identity.Public{
				NodeID: fmt.Sprintf("100 000 %03d", idx),
				UUID:   fmt.Sprintf("client-uuid-%d", idx),
			}
			c := NewClient(cid, 2*time.Second)
			res, err := c.Ping(context.Background(), serverAddr)
			if err != nil {
				errCh <- fmt.Errorf("client %d ping error: %w", idx, err)
				return
			}
			if res.PeerID != serverID.NodeID {
				errCh <- fmt.Errorf("client %d unexpected peer: %s", idx, res.PeerID)
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatal(err)
	}
}

func TestServerCleanStop(t *testing.T) {
	srv := NewServer(identity.Public{NodeID: "555 555 555"}, WithServerPort(0), WithServerListenIP("127.0.0.1"))
	if err := srv.Start(context.Background()); err != nil {
		t.Fatalf("server start: %v", err)
	}
	port := srv.Port()

	// Connect a raw socket to verify it gets closed on Stop
	conn, err := net.Dial("tcp4", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if err := srv.Stop(); err != nil {
		t.Fatalf("server stop: %v", err)
	}

	// Attempting to read from client should now return EOF because server closed connection
	buf := make([]byte, 10)
	_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	_, err = conn.Read(buf)
	if err == nil {
		t.Fatal("expected read error after server stop, got nil")
	}
}
