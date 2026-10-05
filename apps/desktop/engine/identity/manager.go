package identity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	dirName            = "Nodex"
	identityDirName    = "identity"
	keyFileName        = "machine.key"
	keyFilePermissions = 0o600
	dirPermissions     = 0o700

	sealedSeedField = "sealed_private_key"
	bindingField    = "machine_binding"
	keyVersion      = "v1"
)

// stored is the on-disk shape of the machine identity. The private key never sits in
// the file as plaintext: it is sealed with a key derived from the machine binding, so
// copying machine.key to another host yields an unusable blob rather than an identity.
type stored struct {
	Version   string `json:"version"`
	NodeID    string `json:"node_id"`
	UUID      string `json:"uuid"`
	PublicKey string `json:"public_key"`
	SealedKey string `json:"sealed_private_key"`
	Binding   string `json:"machine_binding"`
}

// Public is the subset of the identity that is safe to hand to the frontend.
type Public struct {
	NodeID    string `json:"node_id"`
	UUID      string `json:"uuid"`
	PublicKey string `json:"public_key"`
}

// Identity is the runtime representation of this installation's permanent identity.
type Identity struct {
	NodeID    string
	UUID      string
	PublicKey string
	Seed      ed25519.PrivateKey
}

// Public strips the private key material from the identity.
func (i Identity) Public() Public {
	return Public{NodeID: i.NodeID, UUID: i.UUID, PublicKey: i.PublicKey}
}

// KeyPair is derived cryptographic material for this installation.
type KeyPair struct {
	PublicKey ed25519.PublicKey
	Seed      ed25519.PrivateKey
}

// Manager owns the on-disk lifecycle of the machine identity.
type Manager struct {
	path     string
	binding  string
	identity Identity
}

// NewManager resolves the OS-standard config path for the machine key file.
func NewManager() (*Manager, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("identity: resolve user config dir: %w", err)
	}
	binding, err := MachineBinding()
	if err != nil {
		return nil, err
	}
	return &Manager{
		path:    filepath.Join(base, dirName, identityDirName, keyFileName),
		binding: binding,
	}, nil
}

// Path returns the absolute path of the persisted machine key file.
func (m *Manager) Path() string { return m.path }

// Load reads the persisted identity, generating and persisting a new one on first boot.
func (m *Manager) Load() error {
	raw, err := os.ReadFile(m.path)
	switch {
	case err == nil:
		return m.restore(raw)
	case errors.Is(err, os.ErrNotExist):
		return m.generate()
	default:
		return fmt.Errorf("identity: read %s: %w", m.path, err)
	}
}

// restore decrypts an existing machine key file back into the in-memory identity.
func (m *Manager) restore(raw []byte) error {
	var payload stored
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("identity: parse %s: %w", m.path, err)
	}
	if payload.Version != keyVersion {
		return fmt.Errorf("identity: unsupported key file version %q", payload.Version)
	}
	if payload.Binding != m.binding {
		return ErrBindingMismatch
	}

	seed, err := m.unseal(payload)
	if err != nil {
		return err
	}

	private := ed25519.NewKeyFromSeed(seed)
	m.identity = Identity{
		NodeID:    payload.NodeID,
		UUID:      payload.UUID,
		PublicKey: hex.EncodeToString(private.Public().(ed25519.PublicKey)),
		Seed:      private,
	}
	return nil
}

// generate creates fresh cryptographic material and writes it to disk.
func (m *Manager) generate() error {
	kp, err := GenerateKeyPair()
	if err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("identity: generate uuid: %w", err)
	}

	seed := kp.Seed.Seed()
	nodeID := FormatNodeID(DeriveNodeID(seed))

	payload := stored{
		Version:   keyVersion,
		NodeID:    nodeID,
		UUID:      id.String(),
		PublicKey: hex.EncodeToString(kp.PublicKey),
		Binding:   m.binding,
	}
	sealed, err := m.seal(seed, payload.UUID)
	if err != nil {
		return err
	}
	payload.SealedKey = sealed

	if err := m.write(payload); err != nil {
		return err
	}

	private := ed25519.NewKeyFromSeed(seed)
	m.identity = Identity{
		NodeID:    nodeID,
		UUID:      payload.UUID,
		PublicKey: payload.PublicKey,
		Seed:      private,
	}
	return nil
}

func (m *Manager) write(payload stored) error {
	if err := os.MkdirAll(filepath.Dir(m.path), dirPermissions); err != nil {
		return fmt.Errorf("identity: create identity dir: %w", err)
	}
	encoded, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("identity: encode identity: %w", err)
	}
	if err := os.WriteFile(m.path, encoded, keyFilePermissions); err != nil {
		return fmt.Errorf("identity: write %s: %w", m.path, err)
	}
	return os.Chmod(m.path, keyFilePermissions)
}

// GenerateKeyPair builds an ed25519 keypair from 32 bytes of crypto/rand seed material.
func GenerateKeyPair() (KeyPair, error) {
	seed := make([]byte, ed25519.SeedSize)
	if _, err := rand.Read(seed); err != nil {
		return KeyPair{}, fmt.Errorf("identity: read entropy: %w", err)
	}
	private := ed25519.NewKeyFromSeed(seed)
	return KeyPair{PublicKey: private.Public().(ed25519.PublicKey), Seed: private}, nil
}

// deriveKEK stretches the machine binding into an AES-256 key, salted per identity.
func (m *Manager) deriveKEK(uuid string) ([]byte, error) {
	salt := sha256.Sum256([]byte(uuid))
	key, err := hkdf.Key(sha256.New, []byte(m.binding), salt[:], "nodex identity "+keyVersion, 32)
	if err != nil {
		return nil, fmt.Errorf("identity: derive key: %w", err)
	}
	return key, nil
}

func (m *Manager) seal(seed []byte, uuid string) (string, error) {
	gcm, err := m.cipher(uuid)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("identity: read entropy: %w", err)
	}
	sealed := gcm.Seal(nonce, nonce, seed, []byte(uuid))
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (m *Manager) unseal(payload stored) ([]byte, error) {
	raw, err := base64.StdEncoding.DecodeString(payload.SealedKey)
	if err != nil {
		return nil, fmt.Errorf("identity: decode sealed key: %w", err)
	}

	gcm, err := m.cipher(payload.UUID)
	if err != nil {
		return nil, err
	}
	if len(raw) < gcm.NonceSize() {
		return nil, errors.New("identity: sealed key is truncated")
	}

	seed, err := gcm.Open(nil, raw[:gcm.NonceSize()], raw[gcm.NonceSize():], []byte(payload.UUID))
	if err != nil {
		return nil, fmt.Errorf("identity: unseal private key: %w", err)
	}
	if len(seed) != ed25519.SeedSize {
		return nil, errors.New("identity: unsealed private key has wrong length")
	}
	return seed, nil
}

func (m *Manager) cipher(uuid string) (cipher.AEAD, error) {
	kek, err := m.deriveKEK(uuid)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(kek)
	if err != nil {
		return nil, fmt.Errorf("identity: init cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("identity: init gcm: %w", err)
	}
	return gcm, nil
}
