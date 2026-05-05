# NodeOS Internals & Security Protocols

## Cryptographic Signing
All alerts and logs are signed.
- **Algorithm:** Ed25519
- **Identity:** Generated on-first-boot (Rebirth).
- **Verification:** Can be verified using the node's public key found in `node_id.txt`.

## The Awareness Loop
The node runs a 30-second heart-beat.
1. **Poll Sensors:** Checks for debugger presence and binary changes.
2. **Genesis Check:** Verifies that the Node ID and Public Key match the locked Genesis Hash.
3. **Consult Policy:** Checks autonomy level against manifest.
4. **Action:** If drift is detected, transitions state (Aware -> Lockdown).

## Self-Destruct Sequence
Triggered by extreme manifest violation or physical tampering.
- **Wipe:** Private keys zeroed in memory.
- **Delete:** sensitive files erased from disk.
- **Shutdown:** Process termination.

## Maintenance Protocol
Runs every 24 hours.
- **Time Sync:** Restarts `systemd-timesyncd` to ensure log accuracy.
- **Update:** Fetches only security-related OS patches.
- **SSH Harden:** Re-applies the 22009 port and PubKey-only policy.
