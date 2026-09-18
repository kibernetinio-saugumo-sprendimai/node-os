# NodeOS Lifecycle States

## 1. INIT
The system first checks for a persisted lockdown and refuses normal startup if one exists. It then loads the identity and configures the firewall.

## 2. IDENTITY_BORN
The node has either loaded an existing identity or performed a **REBIRTH**.
- A new Node ID is salted with a hardware identifier when one can be read. Existing identities are not hardware-bound.
- Private Key is loaded into memory.

## 3. AWARE
The system is active and monitoring its environment.
- Ticker runs every 30 seconds.
- Integrity checks are performed.

## 4. DEGRADED
Anomalies detected, but not critical enough for lockdown.
- System functions but alerts frequently.

## 5. LOCKDOWN
Integrity violation detected.
- Network is isolated.
- SSH disabled.
- A persistent lockdown marker prevents systemd restarts from restoring normal firewall rules. Manual intervention is required to clear it.

## 6. TERMINATING / TERMINATED
Self-destruct sequence active or complete.
- Identity wiped.
- Process exited.
- A persistent destruction marker prevents automatic identity recreation on restart; recovery requires an administrator to review and clear it.
