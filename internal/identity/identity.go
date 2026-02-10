package identity

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"pulse/internal/config"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
)

// StoredKey holds the private key on disk (base64-encoded protobuf).
type StoredKey struct {
	PeerID     string `json:"peer_id"`
	PrivateKey string `json:"private_key"`
	Salt       string `json:"salt"`
}

// Generate creates a new Ed25519 identity and saves it to disk.
func Generate() (string, error) {
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return "", fmt.Errorf("generating key: %w", err)
	}

	pid, err := peer.IDFromPrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("deriving peer ID: %w", err)
	}

	raw, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return "", fmt.Errorf("marshaling key: %w", err)
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}

	sk := StoredKey{
		PeerID:     pid.String(),
		PrivateKey: base64.StdEncoding.EncodeToString(raw),
		Salt:       base64.StdEncoding.EncodeToString(salt),
	}

	data, err := json.MarshalIndent(sk, "", "  ")
	if err != nil {
		return "", fmt.Errorf("encoding key: %w", err)
	}

	path := config.IdentityKeyPath()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return "", fmt.Errorf("writing key file: %w", err)
	}

	return pid.String(), nil
}

// LoadPrivateKey reads the private key from disk.
func LoadPrivateKey() (crypto.PrivKey, string, error) {
	path := config.IdentityKeyPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("reading key file: %w", err)
	}

	var sk StoredKey
	if err := json.Unmarshal(data, &sk); err != nil {
		return nil, "", fmt.Errorf("decoding key file: %w", err)
	}

	raw, err := base64.StdEncoding.DecodeString(sk.PrivateKey)
	if err != nil {
		return nil, "", fmt.Errorf("decoding private key: %w", err)
	}

	priv, err := crypto.UnmarshalPrivateKey(raw)
	if err != nil {
		return nil, "", fmt.Errorf("unmarshaling private key: %w", err)
	}

	return priv, sk.PeerID, nil
}

// LoadPublicKeyBytes returns the raw public key bytes for the local identity.
func LoadPublicKeyBytes() ([]byte, error) {
	priv, _, err := LoadPrivateKey()
	if err != nil {
		return nil, err
	}
	return crypto.MarshalPublicKey(priv.GetPublic())
}
