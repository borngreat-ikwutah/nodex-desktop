package identity

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/user"
	"sort"
	"strings"
)

// ErrBindingMismatch is returned when the persisted identity was sealed against a
// different machine binding, which means the host fingerprint changed underneath it.
var ErrBindingMismatch = errors.New("identity: machine binding changed, identity cannot be unsealed")

// MachineBinding returns a stable fingerprint of the host, used as entropy input for
// sealing the private key. It is deliberately NOT the identity: a changed hostname or
// network adapter never mints a new Node ID, it only makes the sealed seed unreadable.
func MachineBinding() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("identity: resolve hostname: %w", err)
	}

	uid := "unknown-user"
	if current, err := user.Current(); err == nil {
		uid = current.Uid
	}

	parts := []string{hostname, uid}
	sort.Strings(parts)

	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:]), nil
}
