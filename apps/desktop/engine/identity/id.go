package identity

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

const (
	minNodeID = 100_000_000
	nodeIDMod = uint64(900_000_000)
)

// DeriveNodeID maps arbitrary key material onto a stable, non-zero 9 digit integer.
// Hashing first keeps the distribution uniform regardless of the seed's shape, and
// guarantees the same seed always yields the same Node ID for the life of the install.
func DeriveNodeID(seed []byte) uint64 {
	sum := sha256.Sum256(seed)
	return uint64(binary.BigEndian.Uint32(sum[:4]))%nodeIDMod + minNodeID
}

// FormatNodeID renders a 9 digit integer as spaced groups: XXX XXX XXX.
func FormatNodeID(value uint64) string {
	digits := fmt.Sprintf("%09d", value)
	return fmt.Sprintf("%s %s %s", digits[:3], digits[3:6], digits[6:])
}
