# Crypto Module

This module provides cryptographic services for encrypting and decrypting sensitive data in the application. It uses AES-256-GCM for authenticated encryption with a 32-byte key.

## Key Components

### Service Interface

The `Service` interface defines the cryptographic operations:

```go
type Service interface {
    Encrypt(plaintext string) (string, error)
    Decrypt(ciphertext string) (string, error)
}
```

### Implementation

The `cryptoService` struct implements the `Service` interface:

```go
type cryptoService struct {
    encryptionKey string
}
```

### Factory Functions

The module provides factory functions for creating the crypto service:

```go
// Create with a specific key
func NewService(encryptionKey string) Service {
    return &cryptoService{
        encryptionKey: encryptionKey,
    }
}

// Create from application config
func NewServiceFromConfig() Service {
    return &cryptoService{
        encryptionKey: config.AppConfig.EncryptionKey,
    }
}
```

## Encryption Process

The `Encrypt` method encrypts plaintext using AES-GCM:

1. Decode the base64 encryption key
2. Create a new AES cipher
3. Create a new GCM cipher mode
4. Generate a random nonce
5. Encrypt the plaintext with the nonce
6. Combine the nonce and ciphertext
7. Encode the result as base64

```go
func (s *cryptoService) Encrypt(plaintext string) (string, error) {
    // Implementation details...
}
```

## Decryption Process

The `Decrypt` method decrypts ciphertext using AES-GCM:

1. Decode the base64 ciphertext
2. Decode the base64 encryption key
3. Create a new AES cipher
4. Create a new GCM cipher mode
5. Extract the nonce from the ciphertext
6. Decrypt the ciphertext with the nonce
7. Return the plaintext

```go
func (s *cryptoService) Decrypt(ciphertext string) (string, error) {
    // Implementation details...
}
```

## Usage

The crypto service is used to encrypt sensitive data before storing it in the database:

```go
// Encrypt API credentials
encryptedCredentials, err := cryptoService.Encrypt(apiCredentialsJSON)
if err != nil {
    return fmt.Errorf("failed to encrypt API credentials: %w", err)
}

// Store encrypted credentials in the database
mapping.APICredentials = &encryptedCredentials
```

And to decrypt sensitive data when retrieving it:

```go
// Decrypt API credentials
if mapping.APICredentials != nil {
    decryptedCredentials, err := cryptoService.Decrypt(*mapping.APICredentials)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt API credentials: %w", err)
    }
    // Use decrypted credentials
}
```

## Security Considerations

- Uses AES-256-GCM, a secure authenticated encryption algorithm
- Generates a unique nonce for each encryption operation
- Requires a 32-byte key (256 bits) for AES-256
- Validates key length and format
- Handles errors securely
- Never logs sensitive data