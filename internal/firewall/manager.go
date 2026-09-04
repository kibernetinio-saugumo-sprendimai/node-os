package firewall

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// Init – pradinė ugniasienės konfigūracija (Safe Default)
func run(args ...string) error {
	if output, err := exec.Command("ufw", args...).CombinedOutput(); err != nil {
		return fmt.Errorf("ufw %v failed: %w: %s", args, err, string(output))
	}
	return nil
}

func Init() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	log.Println("🛡️ Initializing Firewall (UFW)...")

	// 1. Reset ir default taisyklės
	for _, args := range [][]string{
		{"--force", "reset"},
		{"default", "deny", "incoming"},
		{"default", "allow", "outgoing"},
		{"allow", "22009/tcp"},
		{"--force", "enable"},
	} {
		if err := run(args...); err != nil {
			return err
		}
	}

	log.Println("✅ Firewall ACTIVE (Default: Deny Incoming).")
	return nil
}

// ApplyLockdown – aklinas užsidarymas
func ApplyLockdown() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	log.Println("🚨 APPLYING NETWORK LOCKDOWN...")

	// Reset removes stale allow rules before applying the fail-closed policy.
	for _, args := range [][]string{
		{"--force", "reset"},
		{"default", "deny", "incoming"},
		{"default", "deny", "outgoing"},
		{"allow", "out", "443/tcp"},
		{"allow", "out", "53"},
		{"--force", "enable"},
	} {
		if err := run(args...); err != nil {
			return err
		}
	}

	log.Println("🔒 Network LOCKDOWN applied. Only HTTPS/DNS allowed outgoing.")
	return nil
}

// ResetToNormal – grįžimas į įprastą būseną
func ResetToNormal() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	if err := run("default", "allow", "outgoing"); err != nil {
		return err
	}
	if err := run("allow", "22009/tcp"); err != nil {
		return err
	}
	log.Println("🔓 Firewall returned to NORMAL (Allow Outgoing & SSH).")
	return nil
}
