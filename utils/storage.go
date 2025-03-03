package utils

import (
	"errors"
	"os"
	"path/filepath"
)

func GetSecretStoragePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".secret")
}

func StoreEncryptedClientSecret(apiKey string, secretStoragePath string, encryptionKey []byte) error {
	encryptedKey, err := Encrypt(apiKey, encryptionKey)
	if err != nil {
		return err
	}

	return os.WriteFile(secretStoragePath, []byte(encryptedKey), 0600)
}

func GetDecryptedClientSecret(secretStoragePath string, encryptionKey []byte) (string, error) {
	data, err := os.ReadFile(secretStoragePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.New("no client secret found")
		}
		return "", err
	}

	return Decrypt(string(data), encryptionKey)
}

func DeleteClientSecret(secretStoragePath string) error {
	return os.Remove(secretStoragePath)
}

func GetClientIdStoragePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".clientId")
}

func StoreEncryptedClientId(apiKey string, idStoragePath string, encryptionKey []byte) error {
	encryptedKey, err := Encrypt(apiKey, encryptionKey)
	if err != nil {
		return err
	}

	return os.WriteFile(idStoragePath, []byte(encryptedKey), 0600)
}

func GetDecryptedClientId(idStoragePath string, encryptionKey []byte) (string, error) {
	data, err := os.ReadFile(idStoragePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.New("no client id found")
		}
		return "", err
	}

	return Decrypt(string(data), encryptionKey)
}

func DeleteClientId(idStoragePath string) error {
	return os.Remove(idStoragePath)
}
