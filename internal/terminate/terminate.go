package terminate

import (
	"os"
	"time"
)

const SecurityViolationExitCode = 78

func Now(code int) {
	// Give time for final alerts to be sent
	time.Sleep(2 * time.Second)
	os.Exit(code)
}

func SecurityViolation() {
	Now(SecurityViolationExitCode)
}
