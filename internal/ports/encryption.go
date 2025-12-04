package ports

// EncryptionService defines the interface for encryption operations
type EncryptionService interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}
