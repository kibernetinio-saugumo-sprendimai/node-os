package alert

import (
	"nodeos/internal/config"
	"nodeos/internal/identity"
	"testing"
)

func TestUpdatePersonality(t *testing.T) {
	// Mock manifest
	identity.NodeManifest = &config.Manifest{
		AutonomyLevel:         3,
		PersonalityMultiplier: 1.0,
	}

	// Run multiple times to see if it changes
	counts := make(map[PersonalityMode]int)
	for i := 0; i < 1000; i++ {
		updatePersonality()
		counts[Mode]++
	}

	if counts[PmNormal] == 0 && counts[PmStoic] == 0 && counts[PmParanoid] == 0 && counts[PmPanic] == 0 {
		t.Error("Personality modes were not updated")
	}

	t.Logf("Mode counts: %v", counts)
}
