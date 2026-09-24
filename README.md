# SafeStack NodeOS
> **Sovereign, Autonomous, Zero-Trust Defense for Edge Infrastructure.**

SafeStack NodeOS is a specialized security layer designed for mission-critical edge computing devices and autonomous nodes (optimized for Raspberry Pi 5). It transforms standard hardware into a self-auditing, tamper-proof environment through cryptographic anchoring, continuous behavioral awareness, and strict alignment with SafeStack philosophy.

---

## 🏛️ Canonical Alignment & Governance

- **Ecosystem:** SafeStack Decentralized Security Framework
- **Root Canon Compliance:** [`safestack-canon`](https://github.com/kibernetinio-saugumo-sprendimai/safestack-canon) v1.0.0
- **Technical Canon Compliance:** [`safestack-technical-canon`](https://github.com/kibernetinio-saugumo-sprendimai/safestack-technical-canon) v1.0.0
- **Project Identity:** Registered as `project-005` in [`safestack-project-public-keys`](https://github.com/kibernetinio-saugumo-sprendimai/safestack-project-public-keys)
  - Public Key: `A5GWe1zzjq43Rjg9xfNgdSa2u8reF5Z4AkVyvHFoZtw=`
  - Key Fingerprint: `SHA256:a89d861185c7603c0c01722150ee5ea6d61c5da94cabe99973baa41ab9a9c288`
- **Validation Registry:** [`safestack-validation-registry`](https://github.com/kibernetinio-saugumo-sprendimai/safestack-validation-registry)
- **Release Version:** `v0.2.0` (Genesis Core)

---

## 🛡️ Core Pillars

### 1. Hardware-Anchored Identity (Genesis)
Identity is not just a file; it's a bond. NodeOS fingerprints the underlying hardware (SoC Serial / UUID) to generate a unique **Genesis Hash**. Any attempt to move the disk or emulate the hardware results in an immediate identity mismatch and system lockdown.

### 2. Continuous Integrity Verification
The **Awareness Loop** monitors the system every 30 seconds:
- **Binary Integrity**: Verifies the NodeOS executable against its signed hash.
- **Genesis Check**: Ensures the identity anchors haven't been tampered with.
- **Environment Hardening**: Detects debuggers, root privilege escalation, and attack surface expansion (unauthorized open ports).

### 3. Fail-Closed Architecture & Autonomy Severance
Security is binary. In accordance with SafeStack Layer 5 Autonomy Severance:
- **LOCKDOWN**: Instant network isolation (UFW-based drop all inbound/outbound, disable SSH daemons).
- **AUTONOMY SEVERANCE**: Cryptographic erasure of all identity keys and state files, refusing to exist under violation.

### 4. Signed Audit Trails & Manifests
Every system event, alert, and log entry is cryptographically signed using the Node's Ed25519 private key. The security policies in `config/manifest.json` are verified against the embedded Root Key before runtime activation.

---

## 🚀 Quick Start (Pi 5 Optimized)

### Prerequisites
- Raspberry Pi 5 (Recommended: Argon NEO 5 M.2 NVMe case).
- Go 1.22+.
- `smartmontools` (for NVMe health monitoring).

### Installation & Run
```bash
make setup
make build
make run
```

### Verification
Verify project integrity, manifest signatures, and audit trails:
```bash
make verify
```

---

## 📊 Performance & Edge Optimization
NodeOS is engineered for edge deployment:
- **Low Footprint**: Optimized Go binary with stripped debug info (`-ldflags="-s -w"`).
- **Non-Intrusive**: Uses kernel `/proc` inspection instead of active scanning.
- **NVMe Optimized**: High-speed signed logging and SMART health monitoring.

---

## ⚖️ License
Licensed under the SafeStack Ecosystem guidelines. Users maintain 100% sovereignty over their Root Keys.
