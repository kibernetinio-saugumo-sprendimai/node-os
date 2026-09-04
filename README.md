# SafeStack NodeOS
> **Sovereign, Autonomous, Zero-Trust Defense for Edge Infrastructure.**

SafeStack NodeOS is an experimental security layer for OSINT nodes and edge devices (optimized for Raspberry Pi 5). It provides self-auditing, hardware-aware cryptographic anchoring and fail-closed integrity checks. It is not described as tamper-proof and requires target-hardware validation before production use.

---

## 🛡️ Core Pillars

### 1. Hardware-Anchored Identity (Genesis)
Identity is not just a file; it's a bond. NodeOS fingerprints the underlying hardware (SoC Serial / UUID) to generate a unique **Genesis Hash**. Any attempt to move the disk or emulate the hardware results in an immediate identity mismatch and system lockdown.

### 2. Continuous Integrity Verification
The **Awareness Loop** monitors the system every 30 seconds:
- **Binary Integrity**: Verifies the NodeOS executable against its signed hash.
- **Genesis Check**: Ensures the identity anchors haven't been tampered with.
- **Environment Hardening**: Detects debuggers, root privilege escalation, and attack surface expansion (unauthorized open ports).

### 3. Fail-Closed Architecture
Security is binary. If a critical violation is detected, NodeOS responds according to the signed **Manifest**:
- **LOCKDOWN**: Instant network isolation (UFW-based).
- **SELF-DESTRUCT**: Cryptographic erasure of all identity keys and sensitive state files.

### 4. Signed Audit Trails
Every system event, alert, and log entry is cryptographically signed using the Node's Ed25519 private key. This creates an immutable audit trail, verifiable against the user's **Root of Trust**.

---

## 🚀 Quick Start (Pi 5 Optimized)

### Prerequisites
- Raspberry Pi 5 (Recommended: Argon NEO 5 M.2 NVMe case).
- Go 1.22+.
- `smartmontools` (for NVMe health monitoring).

### Installation
```bash
make setup
make build
make run
```

### Verification
Ensure your files are authentic:
```bash
make verify
```

---

## 📊 Performance & Optimization
NodeOS is engineered for the edge:
- **Low Footprint**: Optimized Go binary with stripped debug info.
- **Non-Intrusive**: Uses kernel `/proc` inspection instead of active scanning.
- **NVMe Optimized**: High-speed signed logging and health monitoring.

---

## ⚖️ License & Sovereignty
This project is part of the SafeStack ecosystem. Users maintain 100% sovereignty over their Root Keys.

**SafeStack Security Engineering - Autonomous Systems Division**
