package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stubPasswordReader struct {
	Password    string
	ReturnError bool
}

func (pr stubPasswordReader) ReadPassword() (string, error) {
	if pr.ReturnError {
		return "", errors.New("stubbed error")
	}
	return pr.Password, nil
}

func setupTestStorage() (string, string, func()) {
	tempDir, _ := os.MkdirTemp("", "config_test")
	secretStorage := filepath.Join(tempDir, "test_secret")
	idStorage := filepath.Join(tempDir, "test_client_id")

	cleanup := func() {
		_ = os.Remove(secretStorage)
		_ = os.Remove(idStorage)
		_ = os.RemoveAll(tempDir)
	}

	return secretStorage, idStorage, cleanup
}

func TestConfigStorage(t *testing.T) {
	secretStorage, idStorage, cleanup := setupTestStorage()
	defer cleanup() // Ensure cleanup runs after the test

	// Verify files don't exist before test
	_, err1 := os.Stat(secretStorage)
	_, err2 := os.Stat(idStorage)
	assert.True(t, os.IsNotExist(err1))
	assert.True(t, os.IsNotExist(err2))

	// Create test files
	_ = os.WriteFile(secretStorage, []byte("test-secret"), 0600)
	_ = os.WriteFile(idStorage, []byte("test-client-id"), 0600)

	// Verify files exist after creation
	_, err1 = os.Stat(secretStorage)
	_, err2 = os.Stat(idStorage)
	assert.NoError(t, err1)
	assert.NoError(t, err2)

	// Remove files and verify they no longer exist
	_ = os.Remove(secretStorage)
	_ = os.Remove(idStorage)

	_, err1 = os.Stat(secretStorage)
	_, err2 = os.Stat(idStorage)
	assert.True(t, os.IsNotExist(err1)) // Now correct
	assert.True(t, os.IsNotExist(err2)) // Now correct
}

func TestConfigCommand(t *testing.T) {
	// Simulate user input
	secretStorage, idStorage, cleanup := setupTestStorage()
	defer cleanup()

	input := "test-secret\ntest-client-id\n"
	r, w, _ := os.Pipe()
	_, _ = w.Write([]byte(input))
	_ = w.Close()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	// Capture output
	var output bytes.Buffer
	cmd := NewConfigCmd(secretStorage, idStorage, []byte("1234567891234567"), stubPasswordReader{Password: "password"})
	cmd.SetOut(&output)
	cmd.SetErr(&output)

	// Execute command
	err := cmd.Execute()
	assert.NoError(t, err)
}
