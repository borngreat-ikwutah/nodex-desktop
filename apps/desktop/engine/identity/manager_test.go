package identity

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testBinding = "test-machine-binding"

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	return &Manager{
		path:    filepath.Join(t.TempDir(), dirName, identityDirName, keyFileName),
		binding: testBinding,
	}
}

func TestFirstBootGeneratesIdentity(t *testing.T) {
	m := newTestManager(t)
	if err := m.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}

	if got := len(m.identity.NodeID); got != len("482 913 027") {
		t.Fatalf("node id %q is not a spaced 9 digit id", m.identity.NodeID)
	}
	if m.identity.UUID == "" || m.identity.PublicKey == "" || len(m.identity.Seed) == 0 {
		t.Fatalf("incomplete identity: %+v", m.identity)
	}

	info, err := os.Stat(m.path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != keyFilePermissions {
		t.Fatalf("permissions = %o, want %o", perm, keyFilePermissions)
	}
}

func TestSubsequentBootReusesIdentity(t *testing.T) {
	first := newTestManager(t)
	if err := first.Load(); err != nil {
		t.Fatalf("first load: %v", err)
	}

	second := newTestManager(t)
	second.path = first.path
	if err := second.Load(); err != nil {
		t.Fatalf("second load: %v", err)
	}

	if first.identity.NodeID != second.identity.NodeID {
		t.Fatalf("node id changed across boots: %q vs %q", first.identity.NodeID, second.identity.NodeID)
	}
	if !first.identity.Seed.Equal(second.identity.Seed) {
		t.Fatal("private key changed across boots")
	}
}

func TestPersistedFileNeverContainsPlaintextSeed(t *testing.T) {
	m := newTestManager(t)
	if err := m.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}

	raw, err := os.ReadFile(m.path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var payload stored
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if payload.SealedKey == "" || payload.Version != keyVersion {
		t.Fatalf("unexpected payload: %+v", payload)
	}
	if strings.Contains(string(raw), payload.SealedKey) == false {
		t.Fatal("sealed key missing from payload")
	}
	if payload.NodeID != m.identity.NodeID || payload.UUID != m.identity.UUID {
		t.Fatalf("payload = %+v, want identity %+v", payload, m.identity)
	}

	seed := hex.EncodeToString(m.identity.Seed.Seed())
	if strings.Contains(strings.ToUpper(string(raw)), strings.ToUpper(seed)) {
		t.Fatal("plaintext private key leaked into the key file")
	}
}

func TestCopiedKeyFileIsUnusableOnAnotherMachine(t *testing.T) {
	first := newTestManager(t)
	if err := first.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}

	stolen := newTestManager(t)
	stolen.path = first.path
	stolen.binding = "a-different-machine"

	if err := stolen.Load(); err != ErrBindingMismatch {
		t.Fatalf("err = %v, want ErrBindingMismatch", err)
	}
}

func TestDeriveNodeIDIsStableAndInRange(t *testing.T) {
	seed := []byte("nodex-deterministic-seed")
	first := DeriveNodeID(seed)
	if first != DeriveNodeID(seed) {
		t.Fatal("node id is not deterministic for the same seed")
	}
	if first < minNodeID || first >= minNodeID+nodeIDMod {
		t.Fatalf("node id %d out of range", first)
	}
	if got := FormatNodeID(first); len(got) != 11 || got[3] != ' ' || got[7] != ' ' {
		t.Fatalf("formatted id %q is not XXX XXX XXX", got)
	}
}

func TestLoadRejectsCorruptFile(t *testing.T) {
	m := newTestManager(t)
	if err := os.MkdirAll(filepath.Dir(m.path), dirPermissions); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(m.path, []byte("not json"), keyFilePermissions); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := m.Load(); err == nil {
		t.Fatal("expected error for corrupt identity file")
	}
}

func TestLoadRejectsTamperedSealedKey(t *testing.T) {
	m := newTestManager(t)
	if err := m.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}

	raw, err := os.ReadFile(m.path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	var payload stored
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	payload.SealedKey = payload.SealedKey[:len(payload.SealedKey)-4] + "AAAA"
	tampered, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(m.path, tampered, keyFilePermissions); err != nil {
		t.Fatalf("write: %v", err)
	}

	reloaded := newTestManager(t)
	reloaded.path = m.path
	if err := reloaded.Load(); err == nil {
		t.Fatal("expected error for tampered sealed key")
	}
}

func TestMachineBindingIsStable(t *testing.T) {
	first, err := MachineBinding()
	if err != nil {
		t.Fatalf("binding: %v", err)
	}
	second, err := MachineBinding()
	if err != nil {
		t.Fatalf("binding: %v", err)
	}
	if first != second || first == "" {
		t.Fatalf("binding is not stable: %q vs %q", first, second)
	}
}
