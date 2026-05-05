# NodeOS Security Doctrine

## Zero-Trust Principles
- The host OS is treated as potentially compromised.
- No remote root login.
- No password authentication.
- All outbound communication is restricted during anomalies.

## Threat Response Levels

### Level 1: Warning
Minor drift detected (e.g. non-critical file changes).
- Action: Alert via Telegram. Log signed anomaly.

### Level 2: Lockdown
Significant integrity violation (e.g. binary mismatch, debugger).
- Action: Network isolation via UFW. SSH disabled.

### Level 3: Destruction
Critical security breach (e.g. private key theft attempt, manual manifest violation).
- Action: Wipe keys, wipe logs, delete identity, terminal shutdown.

## Hardening Checklist
- [x] Port 22009 for SSH
- [x] Ed25519 Signed Audit Log
- [x] Hardware Serial Binding
- [x] Immutable Genesis Hash
- [x] Continuous Genesis Integrity Verification
- [x] Anti-Debugging Monitor
