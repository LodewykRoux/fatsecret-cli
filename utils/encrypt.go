package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// var encryptionKey []byte

func LoadEncryptionKey() []byte {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("failed to load .env file")
	}

	key := os.Getenv("ENCRYPTION_KEY")
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		log.Fatal("invalid encryption key length")
	}
	return []byte(key)
}

// Encrypt text using AES
func Encrypt(plaintext string, encryptionKey []byte) (string, error) {
	if plaintext == "" {
		return "", fmt.Errorf("empty text")
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return hex.EncodeToString(ciphertext), nil
}

// Decrypt text using AES
func Decrypt(ciphertextHex string, encryptionKey []byte) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
