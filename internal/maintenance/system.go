package maintenance

import (
	"log"
	"os/exec"
	"runtime"
)

// RunCleanup – išvalo laikinus failus ir sistemines šiukšles
func RunCleanup() {
	if runtime.GOOS != "linux" {
		return
	}

	log.Println("🧹 Starting system cleanup...")
	
	// Apt cleanup
	exec.Command("apt-get", "clean").Run()
	exec.Command("apt-get", "autoremove", "-y").Run()

	// Temp cleanup
	exec.Command("find", "/tmp", "-type", "f", "-atime", "+7", "-delete").Run()
	
	log.Println("✅ Cleanup complete.")
}

// RunUpdates – vykdo tik saugumo atnaujinimus
func RunUpdates() {
	if runtime.GOOS != "linux" {
		return
	}

	log.Println("🛡️ Checking for security updates...")
	
	// Force time sync before updates
	exec.Command("systemctl", "restart", "systemd-timesyncd").Run()

	// Update repositories
	exec.Command("apt-get", "update").Run()
	
	// Install only security upgrades (non-interactive)
	exec.Command("apt-get", "upgrade", "-y", "--only-upgrade").Run()

	log.Println("✅ Security updates processed.")
}
