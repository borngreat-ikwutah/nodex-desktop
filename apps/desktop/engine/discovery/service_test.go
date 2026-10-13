package discovery

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"nodex-desktop/engine/identity"
)

func getFreePort(t *testing.T) int {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer conn.Close()
	return conn.LocalAddr().(*net.UDPAddr).Port
}

func TestNormalizeNodeID(t *testing.T) {
	tests := []struct {
		input   string
		want    string
		wantErr bool
	}{
		{"111222333", "111 222 333", false},
		{"111 222 333", "111 222 333", false},
		{"  111-222-333  ", "111 222 333", false},
		{"12345", "", true},
		{"abcdefghi", "", true},
		{"1234567890", "", true},
	}

	for _, tt := range tests {
		got, err := NormalizeNodeID(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("NormalizeNodeID(%q) expected error, got nil", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("NormalizeNodeID(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("NormalizeNodeID(%q) = %q, want %q", tt.input, got, tt.want)
			}
		}
	}
}

func TestPacketMarshalUnmarshal(t *testing.T) {
	orig := Packet{
		Type:      PacketDiscover,
		TargetID:  "222 222 222",
		SenderID:  "111 111 111",
		SenderIP:  "192.168.1.20",
		UUID:      "mock-uuid-1",
		PublicKey: "mock-pubkey-1",
		Port:      4242,
	}

	data, err := orig.Marshal()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	parsed, err := UnmarshalPacket(data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if parsed.Type != orig.Type || parsed.TargetID != orig.TargetID || parsed.SenderID != orig.SenderID {
		t.Fatalf("mismatch: got %+v, want %+v", parsed, orig)
	}
}

func TestLANDiscoveryBetweenTwoNodes(t *testing.T) {
	testPort := getFreePort(t)

	// Node A: 111 111 111
	idA := identity.Public{
		NodeID:    "111 111 111",
		UUID:      "uuid-node-a",
		PublicKey: "pubkey-node-a",
	}

	// Node B: 222 222 222
	idB := identity.Public{
		NodeID:    "222 222 222",
		UUID:      "uuid-node-b",
		PublicKey: "pubkey-node-b",
	}

	// Spin up Node B as listener on testPort
	serviceB := NewService(idB,
		WithPort(testPort),
		WithListenIP("127.0.0.1"),
		WithTransportPort(4242),
	)

	ctxB, cancelB := context.WithCancel(context.Background())
	defer cancelB()
	if err := serviceB.Startup(ctxB); err != nil {
		t.Fatalf("serviceB startup failed: %v", err)
	}
	defer serviceB.Shutdown(context.Background())

	// Node A sends broadcast targeting testPort
	serviceA := NewService(idA,
		WithPort(testPort),
		WithExtraBroadcastTargets(fmt.Sprintf("127.0.0.1:%d", testPort)),
	)

	findCtx, findCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer findCancel()

	peer, err := serviceA.FindNode(findCtx, "222 222 222")
	if err != nil {
		t.Fatalf("FindNode failed: %v", err)
	}

	if peer.NodeID != "222 222 222" {
		t.Errorf("peer.NodeID = %q, want %q", peer.NodeID, "222 222 222")
	}
	if peer.UUID != "uuid-node-b" {
		t.Errorf("peer.UUID = %q, want %q", peer.UUID, "uuid-node-b")
	}
	if len(peer.IPs) == 0 {
		t.Errorf("expected peer IPs to be non-empty")
	}
}

func TestUnmatchedTargetReturnsNotFound(t *testing.T) {
	testPort := getFreePort(t)

	idB := identity.Public{
		NodeID:    "222 222 222",
		UUID:      "uuid-node-b",
		PublicKey: "pubkey-node-b",
	}
	serviceB := NewService(idB,
		WithPort(testPort),
		WithListenIP("127.0.0.1"),
	)
	if err := serviceB.Startup(context.Background()); err != nil {
		t.Fatalf("startup: %v", err)
	}
	defer serviceB.Shutdown(context.Background())

	idA := identity.Public{
		NodeID: "111 111 111",
		UUID:   "uuid-node-a",
	}
	serviceA := NewService(idA,
		WithPort(testPort),
		WithExtraBroadcastTargets(fmt.Sprintf("127.0.0.1:%d", testPort)),
	)

	// Searching for a non-existent node "333 333 333" should time out
	findCtx, findCancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer findCancel()

	_, err := serviceA.FindNode(findCtx, "333 333 333")
	if !errors.Is(err, ErrNodeNotFound) {
		t.Fatalf("expected ErrNodeNotFound, got: %v", err)
	}
}

func TestSelfDiscoveryRejected(t *testing.T) {
	idA := identity.Public{
		NodeID: "111 111 111",
		UUID:   "uuid-node-a",
	}
	serviceA := NewService(idA)

	_, err := serviceA.FindNode(context.Background(), "111 111 111")
	if err == nil {
		t.Fatal("expected error discovering self, got nil")
	}
}
