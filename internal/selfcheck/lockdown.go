package selfcheck

import (
	"fmt"
	"os"

	"nodeos/internal/alert"
	"nodeos/internal/firewall"
)

const LockdownFlag = "/var/lib/nodeos/LOCKDOWN"

var Locked = false

func EnterLockdown(reason string) error {
	Locked = true
	alert.Critical("NODE LOCKDOWN ACTIVATED: " + reason)
	if err := firewall.ApplyLockdown(); err != nil {
		return fmt.Errorf("apply firewall lockdown: %w", err)
	}

	if err := os.MkdirAll("/var/lib/nodeos", 0700); err != nil {
		return fmt.Errorf("create lockdown directory: %w", err)
	}
	if err := os.WriteFile(LockdownFlag, []byte(reason+"\n"), 0600); err != nil {
		return fmt.Errorf("write lockdown flag: %w", err)
	}

	fmt.Println("=== NODE LOCKDOWN ACTIVATED ===")
	fmt.Println("Reason:", reason)
	fmt.Println("Mode: READ-ONLY")
	fmt.Println("================================")
	return nil
}

func IsLocked() bool {
	if _, err := os.Stat(LockdownFlag); err == nil {
		return true
	}
	return Locked
}
