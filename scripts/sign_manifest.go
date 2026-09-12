package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}

func loadPrivateKey(path string) ed25519.PrivateKey {
	info, err := os.Stat(path)
	if err != nil {
		fail("cannot stat private key file: %v", err)
	}
	if info.Mode().Perm()&0o077 != 0 {
		fail("private key file permissions must be 0600 or stricter")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		fail("cannot read private key file: %v", err)
	}

	keyBytes, err := hex.DecodeString(strings.TrimSpace(string(raw)))
	if err != nil {
		fail("private key must be hex encoded: %v", err)
	}

	switch len(keyBytes) {
	case ed25519.SeedSize:
		return ed25519.NewKeyFromSeed(keyBytes)
	case ed25519.PrivateKeySize:
		return ed25519.PrivateKey(keyBytes)
	default:
		fail("invalid private key length: got %d bytes, want 32 or 64", len(keyBytes))
		return nil
	}
}

func signFile(privateKey ed25519.PrivateKey, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fail("cannot read %s: %v", path, err)
	}

	signature := ed25519.Sign(privateKey, data)
	signaturePath := path + ".sig"
	if err := os.WriteFile(signaturePath, []byte(hex.EncodeToString(signature)), 0o644); err != nil {
		fail("cannot write %s: %v", signaturePath, err)
	}
	fmt.Printf("Signature saved to: %s\n", signaturePath)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run sign_manifest.go <PRIVATE_KEY_FILE> <FILE> [FILE...]")
		os.Exit(2)
	}

	privateKey := loadPrivateKey(os.Args[1])
	for _, path := range os.Args[2:] {
		signFile(privateKey, path)
	}
}
