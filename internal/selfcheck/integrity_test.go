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
	hardwareID := "test-hardware-id"
	pubKey := hex.EncodeToString([]byte("test-public-key-32-bytes-long-now"))
	
	// Reconstruct expected hash
	pubBytes, _ := hex.DecodeString(pubKey)
	genesisData := []byte(id + "|" + hardwareID + "|")
	genesisData = append(genesisData, pubBytes...)
	h := sha256.Sum256(genesisData)
	hashHex := hex.EncodeToString(h[:])

	os.WriteFile("test_genesis.hash", []byte(hashHex), 0644)
	defer os.Remove("test_genesis.hash")

	// 1. Test valid
	sig := CheckGenesisIntegrity(id, pubKey, hardwareID, "test_genesis.hash")
	if sig.Severity != "OK" {
		t.Errorf("Expected OK, got %s", sig.Severity)
	}

	// 2. Test mismatch
	sig = CheckGenesisIntegrity("wrong-id", pubKey, hardwareID, "test_genesis.hash")
	if sig.Severity != "CRITICAL" {
		t.Errorf("Expected CRITICAL for ID mismatch, got %s", sig.Severity)
	}
}

func TestMissingAnchorsFailClosed(t *testing.T) {
	if sig := CheckBinaryIntegrity("does-not-exist"); sig.Severity != "CRITICAL" {
		t.Fatalf("missing binary anchor must be critical, got %s", sig.Severity)
	}
	if sig := CheckGenesisIntegrity("id", hex.EncodeToString(make([]byte, 32)), "hardware", "does-not-exist"); sig.Severity != "CRITICAL" {
		t.Fatalf("missing genesis anchor must be critical, got %s", sig.Severity)
	}
}
