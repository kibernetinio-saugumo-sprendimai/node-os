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

	if identity.GetNodeID() == "" {
		log.Println("No identity found — performing explicit REBIRTH")
		identity.RebirthIdentity()
		alert.Success("REBIRTH ritual complete. New identity generated.")
	}

	log.Println("Node ID:", identity.GetNodeID())

	// -------------------------------
	// MANIFEST (SIGNED)
	// -------------------------------
	if err := identity.LoadManifest(); err != nil {
		log.Println(err)
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

	// Initial check
	decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
	if decision.Lockdown {
		log.Println("SELF-CHECK CRITICAL:", decision.Reason)
		if identity.NodeManifest.SelfDestructOnManifestViolation {
			selfcheck.SelfDestruct(decision.Reason)
		} else if identity.NodeManifest.LockOnTamper {
			selfcheck.EnterLockdown(decision.Reason)
			lc.Transition(lifecycle.Lockdown)
			terminate.Now()
		}
	}

	for {
		select {
		case <-ticker.C:
			decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
			if decision.Lockdown {
				log.Println("SELF-CHECK CRITICAL:", decision.Reason)
				if identity.NodeManifest.SelfDestructOnManifestViolation {
					selfcheck.SelfDestruct(decision.Reason)
				} else if identity.NodeManifest.LockOnTamper {
					selfcheck.EnterLockdown(decision.Reason)
					lc.Transition(lifecycle.Lockdown)
					terminate.Now()
				}
			}
		case <-maintTicker.C:
			maintenance.RunCleanup()
			maintenance.RunUpdates()
			maintenance.CheckNVMeHealth()
		}
	}
}
