# NodeOS Lifecycle States

## 1. INIT
The system is loading core modules and initializing the firewall.

## 2. IDENTITY_BORN
The node has either loaded an existing identity or performed a **REBIRTH**.
- Hardware Serial is checked.
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
- Waiting for manual intervention (Rebirth).

## 6. TERMINATING / TERMINATED
Self-destruct sequence active or complete.
- Identity wiped.
- Process exited.
