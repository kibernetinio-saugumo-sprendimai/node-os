package alert

import (
	"math/rand"
	"nodeos/internal/identity"
)

type PersonalityMode int

const (
	PmNormal PersonalityMode = iota
	PmParanoid
	PmStoic
	PmPanic
)

var Mode PersonalityMode = PmNormal

func updatePersonality() {
	if identity.NodeManifest == nil {
		return
	}

	// Influence by Autonomy Level (0-5)
	// Higher autonomy -> More paranoid/panic chance
	level := float64(identity.NodeManifest.AutonomyLevel)
	multiplier := identity.NodeManifest.PersonalityMultiplier

	x := rand.Float64() * 100 * multiplier

	// Adjust thresholds based on level
	// Base thresholds: 70, 85, 95
	paranoidThreshold := 70.0 - (level * 5)
	stoicThreshold := 85.0 - (level * 2)
	panicThreshold := 95.0 - (level * 1)

	switch {
	case x < paranoidThreshold:
		Mode = PmNormal
	case x < stoicThreshold:
		Mode = PmStoic
	case x < panicThreshold:
		Mode = PmParanoid
	default:
		Mode = PmPanic
	}
}

func personalityComment() string {
	switch Mode {
	case PmParanoid:
		return "\n🧠 pastaba: esu atsargesnis nei įprastai"
	case PmStoic:
		return "\n🧠 pastaba: laikau stabilų toną"
	case PmPanic:
		return "\n🧠 pastaba: reakcijos paaštrintos"
	default:
		return ""
	}
}
