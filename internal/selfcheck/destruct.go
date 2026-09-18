package selfcheck

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"nodeos/internal/alert"
	"nodeos/internal/identity"
)

const DestroyedFlag = "/var/lib/nodeos/SELF_DESTRUCTED"

func IsDestroyed() bool {
	_, err := os.Stat(DestroyedFlag)
	return err == nil
}

func SelfDestruct(reason string) bool {
	alert.Destruct("NODE SELF-DESTRUCT SEQUENCE STARTED: " + reason)
	if err := os.WriteFile(DestroyedFlag, []byte("NODE TERMINATED\n"), 0600); err != nil {
		fmt.Println("CRITICAL: could not persist destruction marker:", err)
		return false
	}
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

	fmt.Println("=======================================")
	fmt.Println(" 🧨 NODE EXECUTED SELF-DESTRUCT")
	fmt.Println(" 🧩 ALL IDENTITY & STATE ERASED")
	fmt.Println(" 🚫 PROCESS TERMINATED")
	fmt.Println("=======================================")

	os.Exit(0)
	return true
}
