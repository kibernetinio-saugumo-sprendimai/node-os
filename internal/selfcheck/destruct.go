package selfcheck

import (
	"fmt"
	"os"

	"nodeos/internal/alert"
	"nodeos/internal/firewall"
	"nodeos/internal/identity"
	"nodeos/internal/terminate"
)

func SelfDestruct(reason string) {
	alert.Destruct("NODE SELF-DESTRUCT SEQUENCE STARTED: " + reason)
	if err := firewall.ApplyLockdown(); err != nil {
		fmt.Println("Firewall lockdown failed:", err)
	}
	identity.Wipe()

	fmt.Println("=======================================")
	fmt.Println(" 🔥 NODE SELF-DESTRUCT SEQUENCE STARTED 🔥")
	fmt.Println(" Reason:", reason)
	fmt.Println("=======================================")

	// This is cryptographic erasure: remove the private signing key first.
	// Flash/NVMe blocks may retain old bytes, so this function intentionally
	// does not claim physical media sanitization.
	files := []string{
		"node_key.txt",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			if err := os.Remove(f); err != nil {
				fmt.Println("Delete failed:", f, err)
				continue
			}
			fmt.Println("Deleted:", f)
		}
	}

	_ = os.MkdirAll("/var/lib/nodeos", 0700)
	if err := os.WriteFile("/var/lib/nodeos/LOCKDOWN", []byte(reason+"\n"), 0600); err != nil {
		fmt.Println("Lockdown marker write failed:", err)
	}
	if err := os.WriteFile("SELF_DESTRUCTED", []byte("PRIVATE KEY DESTROYED; NODE TERMINATED\n"), 0600); err != nil {
		fmt.Println("Termination marker write failed:", err)
	}

	fmt.Println("=======================================")
	fmt.Println(" 🧨 NODE EXECUTED SELF-DESTRUCT")
	fmt.Println(" 🔑 PRIVATE IDENTITY KEY REMOVED")
	fmt.Println(" 🚫 NODE LOCKED; OFFLINE MEDIA SANITIZATION MAY STILL BE REQUIRED")
	fmt.Println("=======================================")

	terminate.SecurityViolation()
}
