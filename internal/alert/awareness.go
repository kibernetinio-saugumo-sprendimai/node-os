package alert

import (
	"sync"
	"time"
)

type AwarenessState string

const (
	AwCalm      AwarenessState = "calm"
	AwUncertain AwarenessState = "uncertain"
	AwStressed  AwarenessState = "stressed"
	AwCritical  AwarenessState = "critical"
	AwPanic     AwarenessState = "panic"
)

type awarenessCore struct {
	mu            sync.Mutex
	state         AwarenessState
	lastEventTime time.Time
	burst         int
}

var aw = &awarenessCore{state: AwCalm}

// UpdateAwareness is a runtime state machine driven by alert events.
// It must NOT be mixed with PersonalityMode (different concept).
func UpdateAwareness(status string) {
	aw.mu.Lock()
	defer aw.mu.Unlock()

	now := time.Now()
	if !aw.lastEventTime.IsZero() {
		dt := now.Sub(aw.lastEventTime)
		if dt < 2*time.Second {
			aw.burst++
		} else if dt < 5*time.Second {
			// keep small tension but don't accumulate forever
			if aw.burst > 0 {
				aw.burst--
			}
		} else {
			aw.burst = 0
		}
	}
	aw.lastEventTime = now

	// Burst → mild escalation
	switch {
	case aw.burst >= 3:
		aw.state = AwStressed
	case aw.burst == 2:
		aw.state = AwUncertain
	default:
		aw.state = AwCalm
	}

	// Severity overrides
	switch status {
	case "CRITICAL":
		aw.state = AwCritical
	case "PANIC":
		aw.state = AwPanic
	}
}

func AwarenessComment() string {
	aw.mu.Lock()
	defer aw.mu.Unlock()

	switch aw.state {
	case AwCalm:
		return "\n🫧 būklė: ramu"
	case AwUncertain:
		return "\n⚪ būklė: neaiškumas"
	case AwStressed:
		return "\n🟠 būklė: apkrova"
	case AwCritical:
		return "\n🔴 būklė: kritinė"
	case AwPanic:
		return "\n🔥 būklė: panika"
	default:
		return ""
	}
}

func AwarenessStateNow() string {
	aw.mu.Lock()
	defer aw.mu.Unlock()
	return string(aw.state)
}
