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
)

// main initializes the system components and starts the continuous awareness loop.
func main() {
	logger.Init()
	if selfcheck.IsDestroyed() {
		log.Println("Persisted destruction marker is active; refusing startup.")
		return
	}
	if selfcheck.IsLocked() {
		if err := firewall.ApplyLockdown(); err != nil {
			log.Println("CRITICAL: could not reapply persisted lockdown:", err)
		}
		log.Println("Persisted lockdown is active; refusing normal startup.")
		return
	}
	if err := identity.Init(); err != nil {
		log.Println("Invalid local identity:", err)
		selfcheck.EnterLockdown("INVALID_LOCAL_IDENTITY")
		return
	}
	if err := firewall.Init(); err != nil {
		log.Println("Firewall initialization failed:", err)
		selfcheck.EnterLockdown("FIREWALL_INITIALIZATION_FAILED")
		return
	}
	alert.Info("SafeStack NodeOS starting...")

	lc := lifecycle.New()

	// -------------------------------
	// IDENTITY
	// -------------------------------
	if identity.GetNodeID() == "" {
		log.Println("No identity found — performing explicit REBIRTH")
		if err := identity.RebirthIdentity(); err != nil {
			log.Println("Identity creation failed:", err)
			selfcheck.EnterLockdown("IDENTITY_CREATION_FAILED")
			return
		}
		alert.Success("REBIRTH ritual complete. New identity generated.")
	}

	log.Println("Node ID:", identity.GetNodeID())

	// -------------------------------
	// MANIFEST (SIGNED)
	// -------------------------------
	if err := identity.LoadManifest(); err != nil {
		log.Println(err)
		lc.Transition(lifecycle.Lockdown)
		selfcheck.EnterLockdown("MANIFEST_VERIFICATION_FAILED")
		return
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
			if !selfcheck.SelfDestruct(decision.Reason) {
				selfcheck.EnterLockdown("DESTRUCTION_MARKER_WRITE_FAILED")
				return
			}
		} else {
			selfcheck.EnterLockdown(decision.Reason)
			lc.Transition(lifecycle.Lockdown)
			return
		}
	}

	for {
		select {
		case <-ticker.C:
			decision := selfcheck.Run(identity.NodeManifest.AutonomyLevel)
			if decision.Lockdown {
				log.Println("SELF-CHECK CRITICAL:", decision.Reason)
				if identity.NodeManifest.SelfDestructOnManifestViolation {
					if !selfcheck.SelfDestruct(decision.Reason) {
						selfcheck.EnterLockdown("DESTRUCTION_MARKER_WRITE_FAILED")
						return
					}
				} else {
					selfcheck.EnterLockdown(decision.Reason)
					lc.Transition(lifecycle.Lockdown)
					return
				}
			}
		case <-maintTicker.C:
			maintenance.RunCleanup()
			maintenance.RunUpdates()
			maintenance.CheckNVMeHealth()
		}
	}
}
