package config

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
)

var embeddedRootPubHex = "9760c594fe7e5638a2a6c351db7503817fb803a43cf5ad8547a08d8b6297ad22"

func getEmbeddedRootKey() (ed25519.PublicKey, error) {
	keyBytes, err := hex.DecodeString(embeddedRootPubHex)
	if err != nil {
		return nil, fmt.Errorf("invalid embedded root key")
	}
	return ed25519.PublicKey(keyBytes), nil
}

func VerifyManifestSignature(manifestPath, sigPath string) error {
	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("manifest read failed: %w", err)
	}

	sigHex, err := os.ReadFile(sigPath)
	if err != nil {
		return fmt.Errorf("manifest signature missing")
	}

	sig, err := hex.DecodeString(string(sigHex))
	if err != nil {
		return fmt.Errorf("manifest signature decode failed (invalid hex)")
	}

	pubKey, err := getEmbeddedRootKey()
	if err != nil {
		return err
	}

	if !ed25519.Verify(pubKey, manifest, sig) {
		return fmt.Errorf("manifest signature INVALID")
	}

	return nil
}
