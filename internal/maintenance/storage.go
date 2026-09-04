package maintenance

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"

	"nodeos/internal/alert"
)

type SmartOutput struct {
	SmartStatus struct {
		Passed bool `json:"passed"`
	} `json:"smart_status"`
	NvmeSmartHealthLog struct {
		PercentageUsed  int `json:"percentage_used"`
		CriticalWarning int `json:"critical_warning"`
		Temperature     int `json:"temperature"`
	} `json:"nvme_smart_health_information_log"`
}

// CheckNVMeHealth – tikrina NVMe disko sveikatą naudojant smartctl
func CheckNVMeHealth() error {
	if runtime.GOOS != "linux" {
		return nil
	}

	// Bandome nuskaityti pirmo NVMe disko informaciją JSON formatu
	out, err := exec.Command("smartctl", "-a", "/dev/nvme0n1", "--json").Output()
	if err != nil {
		return fmt.Errorf("smartctl failed: %w", err)
	}

	var smart SmartOutput
	if err := json.Unmarshal(out, &smart); err != nil {
		return fmt.Errorf("decode smartctl output: %w", err)
	}

	// 1. Kritinis perspėjimas (Hardware level)
	if smart.NvmeSmartHealthLog.CriticalWarning != 0 {
		alert.Critical(fmt.Sprintf("NVMe HARDWARE WARNING! Critical flag: %d", smart.NvmeSmartHealthLog.CriticalWarning))
	}

	// 2. Nusidėvėjimas (Percentage Used)
	if smart.NvmeSmartHealthLog.PercentageUsed > 80 {
		alert.Warn(fmt.Sprintf("NVMe SSD is wearing out: %d%% used", smart.NvmeSmartHealthLog.PercentageUsed))
	}

	// 3. Temperatūra (Pi 5 + Argon NEO 5 atveju svarbu)
	// smartctl temperatūrą grąžina Kelvinais arba Celsijais priklausomai nuo versijos, 
	// bet naujausios --json versijos dažniausiai grąžina Celsijų.
	tempC := smart.NvmeSmartHealthLog.Temperature
	if tempC > 70 {
		alert.Warn(fmt.Sprintf("NVMe SSD temperature is HIGH: %d°C", tempC))
	}

	// 4. Bendras statusas
	if !smart.SmartStatus.Passed {
		alert.Critical("NVMe SMART status: FAILED!")
		return fmt.Errorf("NVMe SMART status failed")
	}
	return nil
}
