package identity

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// getHardwareID – nuskaito unikalią mikrochemos informaciją (Linux / Pi)
func getHardwareID() (string, error) {
	// 1. Standartinis PC (DMI/UUID)
	if data, err := os.ReadFile("/sys/class/dmi/id/product_uuid"); err == nil {
		if id := strings.TrimSpace(string(data)); id != "" {
			return id, nil
		}
	}
	
	// 2. Raspberry Pi (Serial Number)
	if data, err := os.ReadFile("/proc/device-tree/serial-number"); err == nil {
		if id := strings.TrimSpace(string(data)); id != "" {
			return id, nil
		}
	}

	// 3. Fallback į CPU info (senesnės ARM versijos)
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Serial") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					if id := strings.TrimSpace(parts[1]); id != "" {
						return id, nil
					}
				}
			}
		}
	}

	return "", errors.New("stable hardware identifier unavailable")
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
func loadIdentity() {

	// 1) Load Node ID (NO auto-create)
	idData, err := os.ReadFile("node_id.txt")
	if err != nil {
		Identity = NodeIdentity{} // Empty identity
		return
	}
	nodeID := strings.TrimSpace(string(idData))
	if nodeID == "" {
		Identity = NodeIdentity{}
		return
	}

	// 2) Load private key (NO auto-create)
	keyHex, err := os.ReadFile("node_key.txt")
	if err != nil {
		Identity = NodeIdentity{}
		return
	}

	privBytes, err := hex.DecodeString(strings.TrimSpace(string(keyHex)))
	if err != nil || len(privBytes) != ed25519.PrivateKeySize {
		Identity = NodeIdentity{}
		return
	}

	// derive public key
	pub := privBytes[ed25519.PrivateKeySize-ed25519.PublicKeySize:]

	Identity = NodeIdentity{
		NodeID:     nodeID,
		PrivateKey: ed25519.PrivateKey(privBytes),
		PublicKey:  ed25519.PublicKey(pub),
	}
}

// ===============================
// REBIRTH: CREATE NEW ID + KEYS
// ===============================
func RebirthIdentity() error {

	fmt.Println("⚠️  No identity detected — starting REBIRTH ritual...")

	// 1. Create new Node ID (Salted with Hardware ID)
	hwID, err := getHardwareID()
	if err != nil {
		return err
	}
	randomID := make([]byte, 32)
	if _, err := rand.Read(randomID); err != nil {
		return fmt.Errorf("generate node entropy: %w", err)
	}
	
	// Sukuriame mišrų ID: SHA256(HWID + Random)
	hID := sha256.Sum256(append([]byte(hwID), randomID...))
	newID := hex.EncodeToString(hID[:])[:32] // Naudojame pirmus 32 simbolius kaip ID
	
	if err := os.WriteFile("node_id.txt", []byte(newID), 0600); err != nil {
		return fmt.Errorf("write node id: %w", err)
	}

	// 2. Generate new ed25519 keypair
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return fmt.Errorf("generate identity key: %w", err)
	}

	encodedPriv := hex.EncodeToString(priv)
	if err := os.WriteFile("node_key.txt", []byte(encodedPriv), 0600); err != nil {
		return fmt.Errorf("write identity key: %w", err)
	}

	// 3. Capture Binary Hash
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}
	b, err := os.ReadFile(exePath)
	if err != nil {
		return fmt.Errorf("read executable: %w", err)
	}
	h := sha256.Sum256(b)
	if err := os.WriteFile("node_bin.hash", []byte(hex.EncodeToString(h[:])), 0400); err != nil {
		return fmt.Errorf("write binary anchor: %w", err)
	}
	fmt.Println("🧬 Binary hash captured and locked.")

	// 4. Generate Genesis Hash (Anchor)
	genesisData := []byte(newID + "|" + hwID + "|")
	genesisData = append(genesisData, pub...)
	gHash := sha256.Sum256(genesisData)
	
	// Sukuriame failą su tik skaitymo teisėmis (0400)
	if err := os.WriteFile("genesis_hash.txt", []byte(hex.EncodeToString(gHash[:])), 0400); err != nil {
		return fmt.Errorf("write genesis anchor: %w", err)
	}
	
	// Linux specifinis užrakinimas (Immutable)
	if runtime.GOOS == "linux" {
		if err := exec.Command("chattr", "+i", "genesis_hash.txt").Run(); err != nil {
			return fmt.Errorf("lock genesis anchor: %w", err)
		}
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
func Init() {
	loadIdentity()
}

func GetNodeID() string {
	return Identity.NodeID
}

func GetPublicKey() string {
	return hex.EncodeToString(Identity.PublicKey)
}

func GetHardwareID() (string, error) {
	return getHardwareID()
}

func Ready() bool {
	return Identity.NodeID != "" &&
		len(Identity.PrivateKey) == ed25519.PrivateKeySize &&
		len(Identity.PublicKey) == ed25519.PublicKeySize
}

func Sign(data []byte) ([]byte, error) {
	if !Ready() {
		return nil, errors.New("node identity is not initialized")
	}
	return ed25519.Sign(Identity.PrivateKey, data), nil
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
