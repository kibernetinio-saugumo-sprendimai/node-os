package maintenance

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// HardenSSH – perrašo sshd_config su saugiais nustatymais
func HardenSSH() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	configPath := "/etc/ssh/sshd_config"
	log.Println("🛡️ Hardening SSH configuration...")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("read sshd_config: %w", err)
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

	tmp, err := os.CreateTemp("/etc/ssh", "sshd_config.safestack.*")
	if err != nil {
		return fmt.Errorf("create temporary sshd_config: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return fmt.Errorf("chmod temporary sshd_config: %w", err)
	}
	if _, err := tmp.WriteString(strings.Join(newLines, "\n")); err != nil {
		tmp.Close()
		return fmt.Errorf("write temporary sshd_config: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync temporary sshd_config: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temporary sshd_config: %w", err)
	}
	if output, err := exec.Command("sshd", "-t", "-f", tmpPath).CombinedOutput(); err != nil {
		return fmt.Errorf("sshd preflight failed: %w: %s", err, string(output))
	}

	backupPath := configPath + ".safestack.bak"
	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("backup sshd_config: %w", err)
	}
	if err := os.Rename(tmpPath, configPath); err != nil {
		return fmt.Errorf("activate sshd_config: %w", err)
	}

	if output, err := exec.Command("systemctl", "restart", "ssh").CombinedOutput(); err != nil {
		_ = os.WriteFile(configPath, data, 0600)
		_, _ = exec.Command("systemctl", "restart", "ssh").CombinedOutput()
		return fmt.Errorf("restart ssh failed; configuration rolled back: %w: %s", err, string(output))
	}
	log.Println("✅ SSH Hardening COMPLETE. Port 22009 active, Password Auth DISABLED.")
	return nil
}
