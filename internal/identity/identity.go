package identity

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

// getHardwareID – nuskaito unikalią mikrochemos informaciją (Linux / Pi)
func getHardwareID() string {
	// 1. Standartinis PC (DMI/UUID)
	if data, err := os.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
		return strings.TrimSpace(string(data))
	}

	// 2. Raspberry Pi (Serial Number)
	if data, err := os.ReadFile("/proc/device-tree/serial-number"); err == nil {
		return strings.TrimSpace(string(data))
	}

	// 3. Fallback į CPU info (senesnės ARM versijos)
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Serial") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	return "GENERIC-HARDWARE-ID"
}

type NodeIdentity struct {
	NodeID     string
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

var Identity NodeIdentity

// ===============================
// LOAD EXISTING IDENTITY (NO AUTO)
// ===============================
func loadIdentity() error {

	// 1) Load Node ID (NO auto-create)
	idData, err := os.ReadFile("node_id.txt")
	if err != nil {
		Identity = NodeIdentity{} // Empty identity
		if !os.IsNotExist(err) {
			return fmt.Errorf("read node_id.txt: %w", err)
		}
		if _, keyErr := os.Stat("node_key.txt"); keyErr == nil {
			return fmt.Errorf("node_key.txt exists but node_id.txt is missing")
		} else if !os.IsNotExist(keyErr) {
			return fmt.Errorf("inspect node_key.txt: %w", keyErr)
		}
		return nil
	}
	nodeID := string(idData)

	// 2) Load private key (NO auto-create)
	keyHex, err := os.ReadFile("node_key.txt")
	if err != nil {
		Identity = NodeIdentity{}
		return fmt.Errorf("read node_key.txt: %w", err)
	}
	keyInfo, err := os.Stat("node_key.txt")
	if err != nil {
		Identity = NodeIdentity{}
		return fmt.Errorf("inspect node_key.txt: %w", err)
	}
	if keyInfo.Mode().Perm()&0077 != 0 {
		Identity = NodeIdentity{}
		return fmt.Errorf("node_key.txt permissions are too broad: %04o", keyInfo.Mode().Perm())
	}

	privBytes, err := hex.DecodeString(string(keyHex))
	if err != nil {
		Identity = NodeIdentity{}
		return fmt.Errorf("decode node_key.txt: %w", err)
	}
	if len(privBytes) != ed25519.PrivateKeySize {
		Identity = NodeIdentity{}
		return fmt.Errorf("invalid Ed25519 private key length: got %d bytes", len(privBytes))
	}
	derived := ed25519.PrivateKey(privBytes).Public().(ed25519.PublicKey)
	if !ed25519.PublicKey(privBytes[ed25519.SeedSize:]).Equal(derived) {
		Identity = NodeIdentity{}
		return fmt.Errorf("node_key.txt contains an inconsistent Ed25519 key")
	}
	if strings.TrimSpace(nodeID) == "" || strings.TrimSpace(nodeID) != nodeID {
		Identity = NodeIdentity{}
		return fmt.Errorf("invalid node_id.txt")
	}

	// derive public key
	pub := privBytes[ed25519.PrivateKeySize-ed25519.PublicKeySize:]

	Identity = NodeIdentity{
		NodeID:     nodeID,
		PrivateKey: ed25519.PrivateKey(privBytes),
		PublicKey:  ed25519.PublicKey(pub),
	}
	return nil
}

// ===============================
// REBIRTH: CREATE NEW ID + KEYS
// ===============================
func RebirthIdentity() error {

	fmt.Println("⚠️  No identity detected — starting REBIRTH ritual...")

	// 1. Create new Node ID (Salted with Hardware ID)
	hwID := getHardwareID()
	rawID := uuid.New().String()

	// Sukuriame mišrų ID: SHA256(HWID + Random)
	hID := sha256.Sum256([]byte(hwID + rawID))
	newID := hex.EncodeToString(hID[:])[:32] // Naudojame pirmus 32 simbolius kaip ID

	if err := os.WriteFile("node_id.txt", []byte(newID), 0644); err != nil {
		return fmt.Errorf("write node_id.txt: %w", err)
	}

	// 2. Generate new ed25519 keypair
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return fmt.Errorf("generate Ed25519 key: %w", err)
	}

	encodedPriv := hex.EncodeToString(priv)
	if err := os.WriteFile("node_key.txt", []byte(encodedPriv), 0600); err != nil {
		return fmt.Errorf("write node_key.txt: %w", err)
	}

	// 3. Capture Binary Hash
	exePath, err := os.Executable()
	if err == nil {
		if b, err := os.ReadFile(exePath); err == nil {
			h := sha256.Sum256(b)
			if err := os.WriteFile("node_bin.hash", []byte(hex.EncodeToString(h[:])), 0644); err != nil {
				return fmt.Errorf("write node_bin.hash: %w", err)
			}
			fmt.Println("🧬 Binary hash captured and locked.")
		} else {
			return fmt.Errorf("read executable for hash: %w", err)
		}
	} else {
		return fmt.Errorf("locate executable for hash: %w", err)
	}

	// 4. Generate Genesis Hash (Anchor)
	genesisData := append([]byte(newID), pub...)
	gHash := sha256.Sum256(genesisData)

	// Sukuriame failą su tik skaitymo teisėmis (0400)
	if err := os.WriteFile("genesis_hash.txt", []byte(hex.EncodeToString(gHash[:])), 0400); err != nil {
		return fmt.Errorf("write genesis_hash.txt: %w", err)
	}

	// Linux specifinis užrakinimas (Immutable)
	if runtime.GOOS == "linux" {
		exec.Command("chattr", "+i", "genesis_hash.txt").Run()
	}

	fmt.Println("⚓ Genesis Hash anchored and LOCKED (immutable).")

	// 5. Apply in-memory
	Identity = NodeIdentity{
		NodeID:     newID,
		PrivateKey: priv,
		PublicKey:  pub,
	}

	fmt.Println("✨ REBIRTH COMPLETE — NEW IDENTITY CREATED ✨")
	fmt.Println("New Node ID:", newID)
	return nil
}

// ===============================
// PUBLIC ACCESS API
// ===============================
func Init() error {
	return loadIdentity()
}

func GetNodeID() string {
	return Identity.NodeID
}

func GetPublicKey() string {
	return hex.EncodeToString(Identity.PublicKey)
}

func Sign(data []byte) []byte {
	return ed25519.Sign(Identity.PrivateKey, data)
}

// Wipe – saugiai ištrina tapatybę iš atminties
func Wipe() {
	if Identity.PrivateKey != nil {
		for i := range Identity.PrivateKey {
			Identity.PrivateKey[i] = 0
		}
	}
	Identity.NodeID = "TERMINATED"
	Identity.PublicKey = nil
	Identity.PrivateKey = nil
}
