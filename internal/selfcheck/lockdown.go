package selfcheck

import (
	"fmt"
	"os"

	"nodeos/internal/alert"
	"nodeos/internal/firewall"
)

const LockdownFlag = "/var/lib/nodeos/LOCKDOWN"

var Locked = false

func EnterLockdown(reason string) {
	Locked = true
	_ = os.MkdirAll("/var/lib/nodeos", 0700)
	_ = os.WriteFile(LockdownFlag, []byte(reason+"\n"), 0600)
	if err := firewall.ApplyLockdown(); err != nil {
		fmt.Println("CRITICAL: firewall lockdown failed:", err)
	}
	alert.Critical("NODE LOCKDOWN ACTIVATED: " + reason)

	fmt.Println("=== NODE LOCKDOWN ACTIVATED ===")
	fmt.Println("Reason:", reason)
	fmt.Println("Mode: READ-ONLY")
	fmt.Println("================================")
}

func IsLocked() bool {
	if _, err := os.Stat(LockdownFlag); err == nil {
		return true
	}
	return Locked
}
