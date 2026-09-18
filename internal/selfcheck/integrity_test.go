package selfcheck

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"testing"

	"nodeos/internal/policy"
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
	if sig.Severity != policy.SeverityOK {
		t.Errorf("Expected OK, got %v: %s", sig.Severity, sig.Reason)
	}

	// 3. Test tampered hash file
	os.WriteFile("test_bin.hash", []byte("wrong-hash"), 0644)
	sig = CheckBinaryIntegrity("test_bin.hash")
	if sig.Severity != policy.SeverityCritical {
		t.Errorf("Expected CRITICAL for tampered hash, got %v", sig.Severity)
	}
}

func TestCheckGenesisIntegrity(t *testing.T) {
	id := "test-node-id"
	pubBytes := make([]byte, 32)
	copy(pubBytes, "test-public-key")
	pubKey := hex.EncodeToString(pubBytes)

	// Reconstruct expected hash
	decodedPub, _ := hex.DecodeString(pubKey)
	genesisData := append([]byte(id), decodedPub...)
	h := sha256.Sum256(genesisData)
	hashHex := hex.EncodeToString(h[:])

	os.WriteFile("test_genesis.hash", []byte(hashHex), 0644)
	defer os.Remove("test_genesis.hash")

	// 1. Test valid
	sig := CheckGenesisIntegrity(id, pubKey, "test_genesis.hash")
	if sig.Severity != policy.SeverityOK {
		t.Errorf("Expected OK, got %v", sig.Severity)
	}

	// 2. Test mismatch
	sig = CheckGenesisIntegrity("wrong-id", pubKey, "test_genesis.hash")
	if sig.Severity != policy.SeverityCritical {
		t.Errorf("Expected CRITICAL for ID mismatch, got %v", sig.Severity)
	}
}

func TestMissingIntegrityReferencesAreCritical(t *testing.T) {
	if sig := CheckBinaryIntegrity("missing-binary-hash"); sig.Severity != policy.SeverityCritical {
		t.Fatalf("missing binary hash should be critical, got %v", sig.Severity)
	}
	if sig := CheckGenesisIntegrity("node", hex.EncodeToString(make([]byte, 32)), "missing-genesis-hash"); sig.Severity != policy.SeverityCritical {
		t.Fatalf("missing genesis anchor should be critical, got %v", sig.Severity)
	}
}
