package terminate

import (
	"os"
	"time"
)

func Now() {
	// Give time for final alerts to be sent
	time.Sleep(2 * time.Second)
	os.Exit(0)
}
