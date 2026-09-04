// Package main provides the entry point for the SafeStack NodeOS.
// SafeStack NodeOS is an autonomous, zero-trust security layer designed for edge devices.
package main

import (
	"log"
	"time"

	"nodeos/internal/alert"
	"nodeos/internal/firewall"
	"nodeos/internal/identity"
	"nodeos/internal/lifecycle"
	"nodeos/internal/logger"
	"nodeos/internal/maintenance"
	"nodeos/internal/selfcheck"
	"nodeos/internal/terminate"
)

// main initializes the system components and starts the continuous awareness loop.
func main() {
	// -------------------------------
	// IDENTITY
	// -------------------------------
	identity.Init()

	reborn := false
	if identity.GetNodeID() == "" {
		log.Println("No identity found — performing explicit REBIRTH")
		if err := identity.RebirthIdentity(); err != nil {
			log.Println("Identity initialization failed:", err)
			terminate.SecurityViolation()
			return
		}
		reborn = true
	}
	if !identity.Ready() {
		log.Println("Identity is incomplete or invalid")
		terminate.SecurityViolation()
		return
	}

	logger.Init()
	lc := lifecycle.New()
	if err := lc.Transition(lifecycle.IdentityBorn); err != nil {
		log.Println("Lifecycle transition failed:", err)
		terminate.SecurityViolation()
		return
	}

	alert.Info("SafeStack NodeOS starting...")
	if reborn {
		alert.Success("REBIRTH ritual complete. New identity generated.")
	}
	log.Println("Node ID:", identity.GetNodeID())

	// -------------------------------
	// MANIFEST (SIGNED)
	// -------------------------------
	if err := identity.LoadManifest(); err != nil {
		securityStop(lc, err.Error())
		return
	}

	log.Println(
		"Manifest loaded. Autonomy:",
		 identity.NodeManifest.AutonomyLevel,
	)
	if err := firewall.Init(); err != nil {
		securityStop(lc, "FIREWALL_INIT_FAILED: "+err.Error())
		return
	}

	// -------------------------------
	// INITIAL MAINTENANCE
	// -------------------------------
	runMaintenance()

	// -------------------------------
	// CONTINUOUS AWARENESS LOOP
	// -------------------------------
	if err := lc.Transition(lifecycle.Aware); err != nil {
		securityStop(lc, "LIFECYCLE_TRANSITION_FAILED: "+err.Error())
		return
	}
	log.Println("NodeOS state:", lc.State())

	ticker := time.NewTicker(30 * time.Second)
	maintTicker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	defer maintTicker.Stop()

	// Initial check
	decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
	if enforceDecision(lc, decision) {
		return
	}

	for {
		select {
		case <-ticker.C:
			decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
			if enforceDecision(lc, decision) {
				return
			}
		case <-maintTicker.C:
			runMaintenance()
		}
	}
}

func enforceDecision(lc *lifecycle.Machine, decision selfcheck.Decision) bool {
	if !decision.Lockdown {
		return false
	}
	log.Println("SELF-CHECK CRITICAL:", decision.Reason)
	if identity.NodeManifest != nil && identity.NodeManifest.SelfDestructOnManifestViolation {
		selfcheck.SelfDestruct(decision.Reason)
		return true
	}
	securityStop(lc, decision.Reason)
	return true
}

func securityStop(lc *lifecycle.Machine, reason string) {
	log.Println("SECURITY STOP:", reason)
	if lc.State() != lifecycle.Lockdown {
		if err := lc.Transition(lifecycle.Lockdown); err != nil {
			log.Println("Lifecycle lockdown transition failed:", err)
		}
	}
	if err := selfcheck.EnterLockdown(reason); err != nil {
		log.Println("Lockdown enforcement error:", err)
	}
	terminate.SecurityViolation()
}

func runMaintenance() {
	checks := []struct {
		name string
		run  func() error
	}{
		{"cleanup", maintenance.RunCleanup},
		{"security updates", maintenance.RunUpdates},
		{"SSH hardening", maintenance.HardenSSH},
		{"NVMe health", maintenance.CheckNVMeHealth},
	}
	for _, check := range checks {
		if err := check.run(); err != nil {
			log.Printf("Maintenance %s failed: %v", check.name, err)
			alert.Warn("MAINTENANCE_FAILED: " + check.name)
		}
	}
}
