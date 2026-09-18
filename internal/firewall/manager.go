package firewall

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func run(args ...string) error {
	if err := exec.Command("ufw", args...).Run(); err != nil {
		return fmt.Errorf("ufw %v: %w", args, err)
	}
	return nil
}

// Init – pradinė ugniasienės konfigūracija (Safe Default)
func Init() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	log.Println("🛡️ Initializing Firewall (UFW)...")

	// 1. Reset ir default taisyklės
	if err := run("--force", "reset"); err != nil {
		return err
	}
	if err := run("default", "deny", "incoming"); err != nil {
		return err
	}
	if err := run("default", "allow", "outgoing"); err != nil {
		return err
	}

	// 2. Leisti SSH specifiniu prievadu
	if err := run("allow", "22009/tcp"); err != nil {
		return err
	}

	// 3. Įjungti
	if err := run("--force", "enable"); err != nil {
		return err
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

	// Reset removes permissive rules left by a previous normal-mode run.
	if err := run("--force", "reset"); err != nil {
		return err
	}
	if err := run("default", "deny", "incoming"); err != nil {
		return err
	}
	if err := run("default", "deny", "outgoing"); err != nil {
		return err
	}
	if err := run("--force", "enable"); err != nil {
		return err
	}
	// Disable SSH daemons so access is revoked independently of firewall state.
	for _, service := range []string{"ssh", "sshd"} {
		if err := exec.Command("systemctl", "disable", "--now", service).Run(); err != nil {
			log.Printf("could not stop %s service: %v", service, err)
		}
	}

	log.Println("🔒 Network LOCKDOWN applied. Inbound and outbound traffic denied.")
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
