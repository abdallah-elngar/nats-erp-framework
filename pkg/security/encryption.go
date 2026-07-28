package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encryptor يقوم بتشفير البيانات
type Encryptor struct {
	key []byte
}

// NewEncryptor ينشئ مشفراً جديداً
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 { // AES-256
		return nil, errors.New("key must be 32 bytes")
	}
	return &Encryptor{key: key}, nil
}

// Encrypt يقوم بتشفير نص
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	// إنشاء GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// إنشاء nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// تشفير
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// تشفير base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt يقوم بفك تشفير نص
func (e *Encryptor) Decrypt(ciphertext string) (string, error) {
	// فك تشفير base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// الحصول على nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]

	// فك التشفير
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
