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
	logger.Init()
	firewall.Init()
	alert.Info("SafeStack NodeOS starting...")

	lc := lifecycle.New()

	// -------------------------------
	// IDENTITY
	// -------------------------------
	identity.Init()

	// Constitutional rule: established identity is never silently replaced.
	// Missing identity is an unknown trust state and therefore fails closed.
	if identity.GetNodeID() == "" {
		reason := "MISSING_ESTABLISHED_IDENTITY"
		log.Println("SELF-CHECK CRITICAL:", reason)
		selfcheck.EnterLockdown(reason)
		lc.Transition(lifecycle.Lockdown)
		terminate.Now()
	}

	log.Println("Node ID:", identity.GetNodeID())

	// -------------------------------
	// MANIFEST (SIGNED)
	// -------------------------------
	if err := identity.LoadManifest(); err != nil {
		reason := "MANIFEST_VALIDATION_FAILED: " + err.Error()
		log.Println(reason)
		selfcheck.EnterLockdown(reason)
		lc.Transition(lifecycle.Lockdown)
		terminate.Now()
	}

	log.Println(
		"Manifest loaded. Autonomy:",
		identity.NodeManifest.AutonomyLevel,
	)

	// -------------------------------
	// INITIAL MAINTENANCE
	// -------------------------------
	maintenance.RunCleanup()
	maintenance.RunUpdates()
	maintenance.HardenSSH()
	maintenance.CheckNVMeHealth()

	// -------------------------------
	// CONTINUOUS AWARENESS LOOP
	// -------------------------------
	lc.Transition(lifecycle.Aware)
	log.Println("NodeOS state:", lc.State())

	ticker := time.NewTicker(30 * time.Second)
	maintTicker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	defer maintTicker.Stop()

	handleCritical := func(reason string) {
		// Constitutional fail-closed behavior. A manifest may not authorize
		// destructive self-erasure or bypass a critical integrity decision.
		log.Println("SELF-CHECK CRITICAL:", reason)
		selfcheck.EnterLockdown(reason)
		lc.Transition(lifecycle.Lockdown)
		terminate.Now()
	}

	// Initial check
	decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
	if decision.Lockdown {
		handleCritical(decision.Reason)
	}

	for {
		select {
		case <-ticker.C:
			decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
			if decision.Lockdown {
				handleCritical(decision.Reason)
			}
		case <-maintTicker.C:
			maintenance.RunCleanup()
			maintenance.RunUpdates()
			maintenance.CheckNVMeHealth()
		}
	}
}
