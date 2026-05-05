# NodeOS Threat Model

## Asset Inventory
1. **Node Private Key:** Critical. Used for signing logs and alerts.
2. **Genesis Hash:** Critical. Anchors the node's existence.
3. **Audit Log:** High. Forensic trail of system activity.

## Threat Actors

### Remote Attacker
- **Goal:** Gain shell access or intercept traffic.
- **Defense:** Hardened SSH, UFW Deny Incoming, No-Password Auth.

### Physical Intermediary
- **Goal:** Steal the SD card or tamper with the hardware.
- **Defense:** HWID binding. If serial mismatch, node refuses to load identity.

### Malicious Operator
- **Goal:** Modify node behavior via manifest.
- **Defense:** Manifest must be signed with the Root Private Key (kept offline).

## Mitigation Strategies
- **Binary Checksumming:** Prevents execution of modified code.
- **Immutable Flags:** Prevents deletion/modification of genesis anchors.
- **Crypto-Logging:** Prevents "erasing footsteps".
