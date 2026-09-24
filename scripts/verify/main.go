package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
)

const rootPubHex = "9760c594fe7e5638a2a6c351db7503817fb803a43cf5ad8547a08d8b6297ad22"

func verifyFile(filePath, sigPath string, pub ed25519.PublicKey) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", filePath, err)
	}

	sigHex, err := os.ReadFile(sigPath)
	if err != nil {
		return fmt.Errorf("read signature %s: %w", sigPath, err)
	}

	sig, err := hex.DecodeString(string(sigHex))
	if err != nil {
		return fmt.Errorf("decode hex signature: %w", err)
	}

	if !ed25519.Verify(pub, data, sig) {
		return fmt.Errorf("signature verification FAILED for %s", filePath)
	}

	return nil
}

func main() {
	pubBytes, err := hex.DecodeString(rootPubHex)
	if err != nil {
		fmt.Println("Error decoding root public key:", err)
		os.Exit(1)
	}
	pub := ed25519.PublicKey(pubBytes)

	files := [][2]string{
		{"config/manifest.json", "config/manifest.json.sig"},
		{"AUDIT_REPORT.md", "AUDIT_REPORT.md.sig"},
	}

	failed := false
	for _, pair := range files {
		if err := verifyFile(pair[0], pair[1], pub); err != nil {
			fmt.Printf("❌ %v\n", err)
			failed = true
		} else {
			fmt.Printf("✅ %s (Ed25519 signature valid)\n", pair[0])
		}
	}

	if failed {
		os.Exit(2)
	}
	fmt.Println("✨ All cryptographic signatures verified successfully.")
}
