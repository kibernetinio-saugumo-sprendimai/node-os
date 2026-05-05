package alert

import (
	"sync"
	"time"
)

type AlertMemory struct {
	mu          sync.Mutex
	LastLevel   string
	LastMessage string
	LastTime    time.Time
	Count1m     int
}

var Memory = AlertMemory{}

func updateMemory(level string, msg string) {
	Memory.mu.Lock()
	defer Memory.mu.Unlock()

	now := time.Now()

	// jei praėjo daugiau nei minutė — reset
	if now.Sub(Memory.LastTime) > time.Minute {
		Memory.Count1m = 0
	}

	Memory.LastLevel = level
	Memory.LastMessage = msg
	Memory.LastTime = now
	Memory.Count1m++
}
