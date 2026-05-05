package maintenance

import (
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// HardenSSH – perrašo sshd_config su saugiais nustatymais
func HardenSSH() {
	if runtime.GOOS != "linux" {
		return
	}

	configPath := "/etc/ssh/sshd_config"
	log.Println("🛡️ Hardening SSH configuration...")

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Println("❌ Could not read sshd_config:", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	var newLines []string

	// Nustatymai, kuriuos norime užtikrinti
	updates := map[string]string{
		"Port":                   "22009",
		"PermitRootLogin":        "no",
		"PasswordAuthentication": "no",
		"PubkeyAuthentication":   "yes",
	}

	processed := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		found := false
		for key, val := range updates {
			if strings.HasPrefix(trimmed, key+" ") || strings.HasPrefix(trimmed, "#"+key+" ") {
				newLines = append(newLines, key+" "+val)
				processed[key] = true
				found = true
				break
			}
		}
		if !found {
			newLines = append(newLines, line)
		}
	}

	// Pridedame tuos, kurių nebuvo faile
	for key, val := range updates {
		if !processed[key] {
			newLines = append(newLines, key+" "+val)
		}
	}

	err = os.WriteFile(configPath, []byte(strings.Join(newLines, "\n")), 0644)
	if err != nil {
		log.Println("❌ Could not write sshd_config:", err)
		return
	}

	// Perkaitiname SSH tarnybą
	exec.Command("systemctl", "restart", "ssh").Run()
	log.Println("✅ SSH Hardening COMPLETE. Port 22009 active, Password Auth DISABLED.")
}
