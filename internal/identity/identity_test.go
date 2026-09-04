package identity

import (
	"os"
	"testing"
)

func TestIdentityLifeCycle(t *testing.T) {
	// Cleanup any existing files
	os.Remove("node_id.txt")
	os.Remove("node_key.txt")
	os.Remove("node_bin.hash")
	os.Remove("genesis_hash.txt")
	t.Cleanup(func() {
		os.Remove("node_id.txt")
		os.Remove("node_key.txt")
		os.Remove("node_bin.hash")
		os.Remove("genesis_hash.txt")
	})

	// 1. Initial State
	Init()
	if GetNodeID() != "" {
		t.Errorf("Expected empty NodeID, got %s", GetNodeID())
	}

	// 2. Rebirth
	if err := RebirthIdentity(); err != nil {
		t.Skipf("stable hardware identity unavailable in test environment: %v", err)
	}
	id := GetNodeID()
	if id == "" {
		t.Fatal("NodeID should not be empty after Rebirth")
	}

	if len(id) != 32 {
		t.Errorf("Expected 32-char NodeID, got %d", len(id))
	}

	// 3. Reload
	Wipe()
	if GetNodeID() != "TERMINATED" {
		t.Error("Wipe failed to clear ID")
	}

	Init()
	if GetNodeID() != id {
		t.Errorf("Reload failed. Expected %s, got %s", id, GetNodeID())
	}

}

func TestSignFailsWithoutIdentity(t *testing.T) {
	previous := Identity
	t.Cleanup(func() { Identity = previous })
	Identity = NodeIdentity{}
	if _, err := Sign([]byte("test")); err == nil {
		t.Fatal("signing without an initialized identity must fail")
	}
}
