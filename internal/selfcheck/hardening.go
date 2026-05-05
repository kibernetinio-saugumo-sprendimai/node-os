package selfcheck

import (
	"fmt"
	"net"
	"nodeos/internal/policy"
	"os"
	"runtime"
	"strings"
)

func CheckEnvironmentHardening() policy.Signal {
	if runtime.GOOS != "linux" {
		return policy.Signal{Severity: policy.SeverityOK}
	}

	// 1. Anti-Debugging (TracerPid check)
	if isBeingTraced() {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonDebuggerDetected,
			Confidence: 1.0,
		}
	}

	// 2. Privilege Check (Must be root for Firewall/Maintenance)
	if os.Geteuid() != 0 {
		return policy.Signal{
			Severity:   policy.SeverityCritical,
			Reason:     policy.ReasonMissingRoot,
			Confidence: 1.0,
		}
	}

	// 3. Attack Surface Audit (Open Ports)
	// Mazgas turėtų turėti tik minimalų kiekį atvirų prievadų (pvz. tik SSH 22009)
	if ports := getOpenPorts(); len(ports) > 3 {
		return policy.Signal{
			Severity:   policy.SeverityWarn,
			Reason:     policy.ReasonAttackSurfaceHigh,
			Confidence: 0.9, // Padidintas pasitikėjimas dėl tikslesnio skenavimo
		}
	}

	return policy.Signal{Severity: policy.SeverityOK}
}

func isBeingTraced() bool {
	data, err := os.ReadFile("/proc/self/status")
	if err != nil {
		return false
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "TracerPid:") {
			parts := strings.Fields(line)
			if len(parts) > 1 && parts[1] != "0" {
				return true
			}
		}
	}
	return false
}

func getOpenPorts() []int {
	data, err := os.ReadFile("/proc/net/tcp")
	if err != nil {
		// Fallback į lėtą metodą, jei /proc nėra (pvz. testuojant ne Linux)
		return getOpenPortsFallback()
	}

	var openPorts []int
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		
		// Tikriname tik "LISTEN" būseną (0A)
		if parts[3] != "0A" {
			continue
		}

		// local_address: port (hex format)
		addrParts := strings.Split(parts[1], ":")
		if len(addrParts) < 2 {
			continue
		}
		
		var port int
		fmt.Sscanf(addrParts[1], "%X", &port)
		if port > 0 {
			openPorts = append(openPorts, port)
		}
	}
	return openPorts
}

func getOpenPortsFallback() []int {
	var openPorts []int
	// Tikriname tik kritinius prievadus, kad nebūtų per lėta
	criticalPorts := []int{21, 22, 23, 25, 80, 443, 3306, 8080}
	for _, port := range criticalPorts {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			openPorts = append(openPorts, port)
		} else {
			ln.Close()
		}
	}
	return openPorts
}
