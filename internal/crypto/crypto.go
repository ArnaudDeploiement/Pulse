package crypto

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/box"
	"lukechampine.com/blake3"
)

const (
	// NonceSize is the size of NaCl box nonce.
	NonceSize = 24
	// KeySize is the size of NaCl keys.
	KeySize = 32
)

// HashFile computes a BLAKE3 hash of the reader content and returns it as bytes.
func HashFile(r io.Reader) ([]byte, error) {
	h := blake3.New(32, nil)
	if _, err := io.Copy(h, r); err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}
	return h.Sum(nil), nil
}

// HashBytes computes a BLAKE3 hash of a byte slice.
func HashBytes(data []byte) []byte {
	h := blake3.Sum256(data)
	return h[:]
}

// GenerateKeyPair generates a NaCl box keypair.
func GenerateKeyPair() (publicKey, privateKey *[KeySize]byte, err error) {
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("generating NaCl keypair: %w", err)
	}
	return pub, priv, nil
}

// Encrypt encrypts plaintext using the recipient's public key and sender's private key.
func Encrypt(plaintext []byte, recipientPub, senderPriv *[KeySize]byte) ([]byte, error) {
	var nonce [NonceSize]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("generating nonce: %w", err)
	}

	encrypted := box.Seal(nonce[:], plaintext, &nonce, recipientPub, senderPriv)
	return encrypted, nil
}

// Decrypt decrypts ciphertext using the sender's public key and recipient's private key.
func Decrypt(ciphertext []byte, senderPub, recipientPriv *[KeySize]byte) ([]byte, error) {
	if len(ciphertext) < NonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	var nonce [NonceSize]byte
	copy(nonce[:], ciphertext[:NonceSize])

	plaintext, ok := box.Open(nil, ciphertext[NonceSize:], &nonce, senderPub, recipientPriv)
	if !ok {
		return nil, fmt.Errorf("decryption failed: authentication error")
	}
	return plaintext, nil
}
