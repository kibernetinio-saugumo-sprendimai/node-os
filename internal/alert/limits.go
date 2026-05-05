package alert

import (
	"sync"
	"time"
)

var (
	lastSent = map[string]time.Time{}
	limitMu  sync.Mutex
)

// Čia apsibrėžiame minimalų laiko tarpą tarp pranešimų
var limits = map[string]time.Duration{
	"INFO":     1 * time.Second,
	"WARN":     3 * time.Second,
	"CRITICAL": 10 * time.Second,
	"ANOMALY":  5 * time.Second,
	"PANIC":    30 * time.Second,
}

// Ši funkcija pasakys: AR leidžiama siųsti alert dabar?
func canSend(level string) bool {
	limitMu.Lock()
	defer limitMu.Unlock()

	now := time.Now()

	// ar šio lygio alert buvo siųstas anksčiau?
	if t, exists := lastSent[level]; exists {
		// ar nepraėjo pakankamai laiko?
		if now.Sub(t) < limits[level] {
			return false // NEGALIMA siųsti
		}
	}

	// Jei praėjo pakankamai laiko → atnaujiname paskutinį siuntimą
	lastSent[level] = now
	return true
}
