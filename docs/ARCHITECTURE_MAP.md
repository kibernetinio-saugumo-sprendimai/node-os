# NodeOS Architecture Map (v0.2.0-genesis)

## Core Philosophy
The NodeOS is a fail-closed, autonomous security layer designed for edge devices (Pi 4/5).

## System Layers

### 1. Identity Layer (`internal/identity`)
- **HWID Binding:** Unique ID salted with Pi Serial.
- **Genesis Anchor:** Immutable hash locked via `chattr +i`.
- **Manifest Loading:** Enforces signed policy execution.

### 2. Security Loop (`internal/selfcheck`)
- **Awareness Engine:** Continuous monitoring of environment.
- **Integrity Monitor:** Binary and manifest hash verification.
- **Hardening:** Anti-debugging and attack surface control.

### 3. Network Defense (`internal/firewall`)
- **UFW Management:** Automated rule generation.
- **Lockdown Mode:** Drastic network isolation on threat detection.
- **Port 22009:** Standardized secure SSH entry.

### 4. Audit Engine (`internal/logger`)
- **Signed Logs:** Every entry is signed with the node's Ed25519 key.
- **Forensic Logs:** Stored in `/var/lib/nodeos/nodeos.crypt.log`.

### 5. Maintenance Engine (`internal/maintenance`)
- **Security Updates:** Non-interactive `apt` security upgrades.
- **Hygiene:** Automated cleanup of system trash and temp files.

---
*SAFE · SILENT · FAIL-SAFE*
