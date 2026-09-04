package selfcheck

import (
	"nodeos/internal/identity"
	"nodeos/internal/policy"
)

func Run(autonomyLevel int) Decision {
	// 1. Check Autonomy Level
	sig := policy.CheckAutonomy(autonomyLevel)
	if sig.Severity == policy.SeverityCritical {
		return Decision{Lockdown: true, Reason: string(sig.Reason)}
	}

	// 2. Check Environment Hardening (Critical)
	sig = CheckEnvironmentHardening()
	if sig.Severity == policy.SeverityCritical {
		return Decision{Lockdown: true, Reason: string(sig.Reason)}
	}

	// 3. Check Binary Integrity (Warn/Critical)
	sig = CheckBinaryIntegrity("node_bin.hash")
	if sig.Severity == policy.SeverityCritical {
		return Decision{Lockdown: true, Reason: string(sig.Reason)}
	}

	// 4. Check Genesis Integrity (Critical)
	// Anchors the Node ID and Public Key to the filesystem
	hardwareID, err := identity.GetHardwareID()
	if err != nil {
		return Decision{Lockdown: true, Reason: "HARDWARE_ID_UNAVAILABLE"}
	}
	sig = CheckGenesisIntegrity(
		identity.GetNodeID(),
		identity.GetPublicKey(),
		hardwareID,
		"genesis_hash.txt",
	)
	if sig.Severity == policy.SeverityCritical {
		return Decision{Lockdown: true, Reason: string(sig.Reason)}
	}

	return Decision{Lockdown: false}
}

type Decision struct {
	Lockdown bool
	Reason   string
}
