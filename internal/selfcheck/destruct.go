package selfcheck

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"nodeos/internal/alert"
	"nodeos/internal/identity"
)

func SelfDestruct(reason string) {
	alert.Destruct("NODE SELF-DESTRUCT SEQUENCE STARTED: " + reason)
	identity.Wipe()

	fmt.Println("=======================================")
	fmt.Println(" 🔥 NODE SELF-DESTRUCT SEQUENCE STARTED 🔥")
	fmt.Println(" Reason:", reason)
	fmt.Println("=======================================")

	files := []string{
		"node_id.txt",
		"node_key.txt",
		"node_bin.hash",
		"genesis_hash.txt",
		"config/nodeos_config.json",
		"config/manifest.json",
		"/var/lib/nodeos/LOCKDOWN",
		"/var/lib/nodeos/nodeos.crypt.log",
	}

	for _, f := range files {
		// Prieš trinant, jei tai genesis_hash, nuimame immutable flag
		if f == "genesis_hash.txt" && runtime.GOOS == "linux" {
			exec.Command("chattr", "-i", f).Run()
		}
		
		if _, err := os.Stat(f); err == nil {
			os.Remove(f)
			fmt.Println("Deleted:", f)
		}
	}

	os.WriteFile("SELF_DESTRUCTED", []byte("NODE TERMINATED\n"), 0644)

	fmt.Println("=======================================")
	fmt.Println(" 🧨 NODE EXECUTED SELF-DESTRUCT")
	fmt.Println(" 🧩 ALL IDENTITY & STATE ERASED")
	fmt.Println(" 🚫 SYSTEM SHUTDOWN")
	fmt.Println("=======================================")

	os.Exit(0)
}
