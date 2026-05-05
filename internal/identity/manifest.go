package identity

import (
	"fmt"
	"os"

	"nodeos/internal/config"
)

var NodeManifest *config.Manifest

func LoadManifest() error {
	manifestPath := os.Getenv("NODEOS_MANIFEST")
	if manifestPath == "" {
		// Try local first, then system
		if _, err := os.Stat("config/manifest.json"); err == nil {
			manifestPath = "config/manifest.json"
		} else {
			manifestPath = "/etc/nodeos/manifest.json"
		}
	}

	sigPath := manifestPath + ".sig"

	// 🔐 1) VERIFY SIGNATURE
	if err := config.VerifyManifestSignature(
		manifestPath,
		sigPath,
	); err != nil {
		return fmt.Errorf("SECURITY VIOLATION: %v", err)
	}

	// 📄 2) LOAD MANIFEST
	m, err := config.LoadManifest(manifestPath)
	if err != nil {
		return fmt.Errorf("manifest load failed: %v", err)
	}

	NodeManifest = m
	return nil
}
