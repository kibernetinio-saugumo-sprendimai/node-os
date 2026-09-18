package selfcheck

import (
	"crypto/sha256"
	"encoding/hex"
	"nodeos/internal/policy"
	"os"
	"strings"
)

func fileSHA256(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

func CheckBinaryIntegrity(hashPath string) policy.Signal {
	exePath, err := os.Executable()
	if err != nil {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonBinaryReadFailed,
			Confidence: 1.0,
		}
	}

	current, err := fileSHA256(exePath)
	if err != nil {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonBinaryReadFailed,
			Confidence: 1.0,
		}
	}

	storedBytes, err := os.ReadFile(hashPath)
	if err != nil {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonBinaryReadFailed,
			Confidence: 0.8,
		}
	}

	stored := strings.TrimSpace(string(storedBytes))
	if stored != current {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonBinaryTampered,
			Confidence: 1.0,
		}
	}

	return policy.Signal{
		Severity:   policy.SeverityOK,
		Reason:     "BINARY_INTEGRITY_OK",
		Confidence: 1.0,
	}
}

func CheckGenesisIntegrity(id string, pubKeyHex string, anchorPath string) policy.Signal {
	storedBytes, err := os.ReadFile(anchorPath)
	if err != nil {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     "GENESIS_ANCHOR_MISSING",
			Confidence: 0.8,
		}
	}

	// Reconstruct expected hash
	pubBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubBytes) != 32 {
		return policy.Signal{Severity: policy.SeverityCritical, Reason: policy.ReasonGenesisMismatch, Confidence: 1.0}
	}
	genesisData := append([]byte(id), pubBytes...)
	h := sha256.Sum256(genesisData)
	current := hex.EncodeToString(h[:])

	stored := strings.TrimSpace(string(storedBytes))
	if stored != current {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonGenesisMismatch,
			Confidence: 1.0,
		}
	}

	return policy.Signal{
		Severity:   policy.SeverityOK,
		Reason:     "GENESIS_INTEGRITY_OK",
		Confidence: 1.0,
	}
}
