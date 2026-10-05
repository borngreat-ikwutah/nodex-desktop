package identity

// Service is the Wails-bound facade over the identity Manager.
type Service struct {
	manager *Manager
}

// NewService loads the machine identity into memory, ready for the frontend.
func NewService() (*Service, error) {
	manager, err := NewManager()
	if err != nil {
		return nil, err
	}
	if err := manager.Load(); err != nil {
		return nil, err
	}
	return &Service{manager: manager}, nil
}

// GetNodeID returns the permanent, human readable Nodex ID of this installation.
func (s *Service) GetNodeID() string { return s.manager.identity.NodeID }

// GetIdentity returns the public identity payload for this installation,
// with private key material stripped out.
func (s *Service) GetIdentity() Public { return s.manager.identity.Public() }

// GetPublicKey returns the hex encoded ed25519 public key of this installation.
func (s *Service) GetPublicKey() string { return s.manager.identity.PublicKey }

// KeyPath returns the absolute path of the persisted machine key file.
func (s *Service) KeyPath() string { return s.manager.Path() }
