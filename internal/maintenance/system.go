package maintenance

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

// RunCleanup – išvalo laikinus failus ir sistemines šiukšles
func runCommand(name string, args ...string) error {
	if output, err := exec.Command(name, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("%s %v failed: %w: %s", name, args, err, string(output))
	}
	return nil
}

func RunCleanup() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	log.Println("🧹 Starting system cleanup...")
	
	// Apt cleanup
	if err := runCommand("apt-get", "clean"); err != nil {
		return err
	}
	if err := runCommand("apt-get", "autoremove", "-y"); err != nil {
		return err
	}

	// Temp cleanup
	if err := runCommand("find", "/tmp", "-xdev", "-type", "f", "-atime", "+7", "-delete"); err != nil {
		return err
	}
	
	log.Println("✅ Cleanup complete.")
	return nil
}

// RunUpdates – vykdo tik saugumo atnaujinimus
func RunUpdates() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	log.Println("🛡️ Checking for security updates...")
	
	// Force time sync before updates
	if err := runCommand("systemctl", "restart", "systemd-timesyncd"); err != nil {
		return err
	}

	// Update repositories
	if err := runCommand("apt-get", "update"); err != nil {
		return err
	}
	
	// Apply only origins allowed by unattended-upgrades policy (normally security).
	if err := runCommand("unattended-upgrade", "-d"); err != nil {
		return err
	}

	log.Println("✅ Security updates processed.")
	return nil
}
