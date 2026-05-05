package config

import (
	"encoding/json"
	"os"
)

type Manifest struct {
	AllowForwardAnomaly  bool    `json:"allow_forward_anomaly"`
	AllowForwardCritical bool    `json:"allow_forward_critical"`
	AllowForwardWarn     bool    `json:"allow_forward_warn"`
	AllowForwardInfo     bool    `json:"allow_forward_info"`

	AutonomyLevel int     `json:"autonomy_level"`

	LockOnTamper                   bool    `json:"lock_on_tamper"`
	SelfDestructOnManifestViolation bool    `json:"self_destruct_on_manifest_violation"`

	PersonalityMultiplier float64 `json:"personality_multiplier"`
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
