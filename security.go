package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var (
	encryptionKey    []byte
	encryptionKeyMux sync.RWMutex
	secretFile       = ".secret"
)

// InitSecurity initializes the encryption key.
// It checks for RTSP_ENCRYPTION_SECRET env var,
// otherwise reads/creates a key in .secret file.
func InitSecurity() {
	encryptionKeyMux.Lock()
	defer encryptionKeyMux.Unlock()

	// 1. Try environment variable
	if envKey := os.Getenv("RTSP_ENCRYPTION_SECRET"); envKey != "" {
		key, err := hex.DecodeString(envKey)
		if err == nil && len(key) == 32 {
			encryptionKey = key
			return
		}
	}

	// 2. Try .secret file
	if _, err := os.Stat(secretFile); err == nil {
		data, err := ioutil.ReadFile(secretFile)
		if err == nil {
			key, err := hex.DecodeString(string(data))
			if err == nil && len(key) == 32 {
				encryptionKey = key
				return
			}
		}
	}

	// 3. Generate new key
	newKey := make([]byte, 32) // AES-256
	if _, err := rand.Read(newKey); err != nil {
		log.Fatalln("Failed to generate encryption key:", err)
	}

	encryptionKey = newKey

	// Save to .secret
	err := ioutil.WriteFile(secretFile, []byte(hex.EncodeToString(newKey)), 0600)
	if err != nil {
		log.Println("Warning: Failed to save encryption key to .secret:", err)
	}
}

// Encrypt string using AES-GCM
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	if strings.HasPrefix(plaintext, "enc:") {
		return plaintext, nil
	}

	encryptionKeyMux.RLock()
	key := encryptionKey
	encryptionKeyMux.RUnlock()

	if len(key) == 0 {
		return plaintext, nil // Fallback if no key (shouldn't happen after Init)
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
	return "enc:" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt string using AES-GCM
func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	if !strings.HasPrefix(ciphertext, "enc:") {
		return ciphertext, nil
	}

	encryptionKeyMux.RLock()
	key := encryptionKey
	encryptionKeyMux.RUnlock()

	if len(key) == 0 {
		return ciphertext, nil
	}

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(ciphertext, "enc:"))
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

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, encryptedData := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
