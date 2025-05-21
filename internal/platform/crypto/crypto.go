package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/affiliate-backend/internal/config"
)

// Service defines the interface for crypto operations
type Service interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

// cryptoService implements the Service interface
type cryptoService struct {
	encryptionKey string
}

// NewService creates a new crypto service
func NewService(encryptionKey string) Service {
	return &cryptoService{
		encryptionKey: encryptionKey,
	}
}

// NewServiceFromConfig creates a new crypto service from the application config
func NewServiceFromConfig() Service {
	return &cryptoService{
		encryptionKey: config.AppConfig.EncryptionKey,
	}
}

// Encrypt encrypts plaintext using AES-GCM with the application's encryption key.
// Returns base64 encoded ciphertext.
func (s *cryptoService) Encrypt(plaintext string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode encryption key: %w", err)
	}
	if len(key) != 32 { // AES-256
		return "", errors.New("encryption key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64 encoded ciphertext using AES-GCM.
func (s *cryptoService) Decrypt(ciphertextB64 string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(s.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode encryption key: %w", err)
	}
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes for AES-256")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, actualCiphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", err // Decryption failed (e.g. wrong key, tampered ciphertext)
	}

	return string(plaintext), nil
}

// For backward compatibility, keep the package-level functions
// that use the application config directly

// Encrypt encrypts plaintext using AES-GCM with the application's encryption key.
// Returns base64 encoded ciphertext.
func Encrypt(plaintext string) (string, error) {
	service := NewServiceFromConfig()
	return service.Encrypt(plaintext)
}

// Decrypt decrypts base64 encoded ciphertext using AES-GCM.
func Decrypt(ciphertextB64 string) (string, error) {
	service := NewServiceFromConfig()
	return service.Decrypt(ciphertextB64)
}