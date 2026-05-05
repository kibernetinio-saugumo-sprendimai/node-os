package config

import (
	"crypto/ed25519"
	"encoding/hex"
	"os"
	"testing"
)

func TestVerifyManifestSignature(t *testing.T) {
	// 1. Setup mock manifest and keys
	pub, priv, _ := ed25519.GenerateKey(nil)
	manifestData := []byte(`{"test": true}`)
	sig := ed25519.Sign(priv, manifestData)
	
	os.WriteFile("test_manifest.json", manifestData, 0644)
	os.WriteFile("test_manifest.json.sig", []byte(hex.EncodeToString(sig)), 0644)
	defer os.Remove("test_manifest.json")
	defer os.Remove("test_manifest.json.sig")

	// Temporarily override root key
	originalKey := embeddedRootPubHex
	embeddedRootPubHex = hex.EncodeToString(pub)
	defer func() { embeddedRootPubHex = originalKey }()

	// 2. Test valid signature
	err := VerifyManifestSignature("test_manifest.json", "test_manifest.json.sig")
	if err != nil {
		t.Errorf("Valid signature failed: %v", err)
	}

	// 3. Test invalid signature
	os.WriteFile("test_manifest.json.sig", []byte("bad-sig"), 0644)
	err = VerifyManifestSignature("test_manifest.json", "test_manifest.json.sig")
	if err == nil {
		t.Error("Expected error for invalid hex signature")
	}

	os.WriteFile("test_manifest.json.sig", []byte(hex.EncodeToString([]byte("wrong-signature-length-64-bytes-needed-here-and-this-is-not-it"))), 0644)
	err = VerifyManifestSignature("test_manifest.json", "test_manifest.json.sig")
	if err == nil {
		t.Error("Expected error for wrong signature")
	}
}
