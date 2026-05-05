package firewall

import (
	"log"
	"os/exec"
	"runtime"
)

// Init – pradinė ugniasienės konfigūracija (Safe Default)
func Init() {
	if runtime.GOOS != "linux" {
		return
	}

	log.Println("🛡️ Initializing Firewall (UFW)...")

	// 1. Reset ir default taisyklės
	exec.Command("ufw", "--force", "reset").Run()
	exec.Command("ufw", "default", "deny", "incoming").Run()
	exec.Command("ufw", "default", "allow", "outgoing").Run()

	// 2. Leisti SSH specifiniu prievadu
	exec.Command("ufw", "allow", "22009/tcp").Run()

	// 3. Įjungti
	exec.Command("ufw", "--force", "enable").Run()
	log.Println("✅ Firewall ACTIVE (Default: Deny Incoming).")
}

// ApplyLockdown – aklinas užsidarymas
func ApplyLockdown() {
	if runtime.GOOS != "linux" {
		return
	}

	log.Println("🚨 APPLYING NETWORK LOCKDOWN...")

	// Blokuojame viską, išskyrus HTTPS (443) pranešimams išsiųsti
	exec.Command("ufw", "default", "deny", "outgoing").Run()
	exec.Command("ufw", "deny", "22009/tcp").Run() // Uždrausti SSH prisijungimą
	exec.Command("ufw", "allow", "out", "443/tcp").Run()
	exec.Command("ufw", "allow", "out", "53").Run() // DNS

	log.Println("🔒 Network LOCKDOWN applied. Only HTTPS/DNS allowed outgoing.")
}

// ResetToNormal – grįžimas į įprastą būseną
func ResetToNormal() {
	if runtime.GOOS != "linux" {
		return
	}

	exec.Command("ufw", "default", "allow", "outgoing").Run()
	exec.Command("ufw", "allow", "22009/tcp").Run()
	log.Println("🔓 Firewall returned to NORMAL (Allow Outgoing & SSH).")
}
