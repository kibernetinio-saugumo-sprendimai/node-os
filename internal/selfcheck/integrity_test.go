package selfcheck

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"
)

func TestCheckBinaryIntegrity(t *testing.T) {
	// 1. Setup mock binary hash
	exePath, _ := os.Executable()
	data, _ := os.ReadFile(exePath)
	h := sha256.Sum256(data)
	hashHex := hex.EncodeToString(h[:])
	
	os.WriteFile("test_bin.hash", []byte(hashHex), 0644)
	defer os.Remove("test_bin.hash")

	// 2. Test valid integrity
	sig := CheckBinaryIntegrity("test_bin.hash")
	if sig.Severity != "OK" {
		t.Errorf("Expected OK, got %s: %s", sig.Severity, sig.Reason)
	}

	// 3. Test tampered hash file
	os.WriteFile("test_bin.hash", []byte("wrong-hash"), 0644)
	sig = CheckBinaryIntegrity("test_bin.hash")
	if sig.Severity != "CRITICAL" {
		t.Errorf("Expected CRITICAL for tampered hash, got %s", sig.Severity)
	}
}

func TestCheckGenesisIntegrity(t *testing.T) {
	id := "test-node-id"
	pubKey := hex.EncodeToString([]byte("test-public-key-32-bytes-long-now"))
	
	// Reconstruct expected hash
	pubBytes, _ := hex.DecodeString(pubKey)
	genesisData := append([]byte(id), pubBytes...)
	h := sha256.Sum256(genesisData)
	hashHex := hex.EncodeToString(h[:])

	os.WriteFile("test_genesis.hash", []byte(hashHex), 0644)
	defer os.Remove("test_genesis.hash")

	// 1. Test valid
	sig := CheckGenesisIntegrity(id, pubKey, "test_genesis.hash")
	if sig.Severity != "OK" {
		t.Errorf("Expected OK, got %s", sig.Severity)
	}

	// 2. Test mismatch
	sig = CheckGenesisIntegrity("wrong-id", pubKey, "test_genesis.hash")
	if sig.Severity != "CRITICAL" {
		t.Errorf("Expected CRITICAL for ID mismatch, got %s", sig.Severity)
	}
}
