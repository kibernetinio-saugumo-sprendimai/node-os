package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run sign_manifest.go <PRIVATE_KEY_HEX> <MANIFEST_FILE>")
		return
	}

	privHex := os.Args[1]
	manifestPath := os.Args[2]

	privBytes, err := hex.DecodeString(privHex)
	if err != nil {
		fmt.Println("Error decoding private key:", err)
		return
	}

	var priv ed25519.PrivateKey
	if len(privBytes) == 32 {
		priv = ed25519.NewKeyFromSeed(privBytes)
	} else if len(privBytes) == 64 {
		priv = ed25519.PrivateKey(privBytes)
	} else {
		fmt.Printf("Error: Invalid private key length (%d bytes). Must be 32 or 64.\n", len(privBytes))
		return
	}

	manifest, err := os.ReadFile(manifestPath)
	if err != nil {
		fmt.Println("Error reading manifest:", err)
		return
	}

	sig := ed25519.Sign(priv, manifest)
	sigHex := hex.EncodeToString(sig)

	fmt.Println("=== MANIFEST SIGNATURE (HEX) ===")
	fmt.Println(sigHex)
	fmt.Println("================================")
	
	sigFile := manifestPath + ".sig"
	err = os.WriteFile(sigFile, []byte(sigHex), 0644)
	if err == nil {
		fmt.Printf("Signature saved to: %s\n", sigFile)
	}
}
