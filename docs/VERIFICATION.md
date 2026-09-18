# Integrity Verification Protocol
> **How to verify the sovereignty and integrity of SafeStack NodeOS.**

This document provides technical instructions for verifying the digital signatures and checksums associated with this release. Following this protocol ensures that the codebase and its audit trails have not been tampered with since the last authorized sign-off.

---

## 🔑 Root of Trust
All signatures are generated using the project's **Root Private Key**. The corresponding public key is hard-coded into the system for runtime verification.

**Root Public Key (Ed25519):**
`9760c594fe7e5638a2a6c351db7503817fb803a43cf5ad8547a08d8b6297ad22`

---

## 🛡️ Step 1: Verify Project Checksums
The `SHA256SUMS` file contains the hashes of every file in the repository. Its integrity is guaranteed by `SHA256SUMS.sig`.

### 1.1 Verify the Signature
To verify that the checksum list itself is authentic, use the following logic (or the provided `make verify` command):
```bash
# Using the built-in Makefile command
make verify
```

### 1.2 Manual Hash Validation
To manually verify that each file matches the registered hash:
```bash
sha256sum -c SHA256SUMS
```

---

## 📜 Step 2: Verify the Audit Report
The `AUDIT_REPORT.md` provides a summary of the security state. Its signature `AUDIT_REPORT.md.sig` ensures the auditor's accountability.

Verification logic:
1. Read the content of `AUDIT_REPORT.md`.
2. Verify it against the hex signature in `AUDIT_REPORT.md.sig` using the **Root Public Key**.

---

## 🧬 Step 3: Verify Runtime Genesis Consistency
When the Node is deployed, it writes a local consistency hash to `genesis_hash.txt`.

The Genesis Hash is derived as follows:
`SHA256(NodeID + NodePublicKey)`

If the file exists and matches this calculation, the Node ID and public key match the values used when this local anchor was written. This is not a hardware binding or an independent trust anchor: an attacker able to replace both the identity files and this hash can create a consistent replacement. Hardware binding requires a separately protected source such as a TPM-backed key or signed provisioning record.

---

## 🖋️ Step 4: Verify the Manifest
The `config/manifest.json` defines the node's behavior. Its signature `manifest.json.sig` prevents unauthorized changes to security policies (like disabling Lockdown).

To verify the manifest:
- The node performs this check **every time it starts**.
- If the manifest signature is invalid, the node will refuse to enter the `Aware` state and will terminate immediately.

---

## 🛠️ Verification Tools
While standard tools like `sha256sum` can verify hashes, Ed25519 signatures require specific utilities. You can use the scripts provided in the `scripts/` directory or any standard Ed25519 library (e.g., `openssl` with provider support or Go's `crypto/ed25519`).

**SafeStack Security Engineering - Verification & Compliance**
