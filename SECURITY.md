# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| v0.2.x  | ✅ YES (Current)    |
| < v0.2  | ❌ NO               |

## Our Security Philosophy

SafeStack NodeOS is built on the principle of **Defensive Autonomy**. We assume the environment is hostile and the hardware might be physically accessed. Our goal is to ensure that even if the node is compromised, it fails closed, protects its identity keys, and provides a signed audit trail of the breach.

## Reporting a Vulnerability

If you discover a security vulnerability within SafeStack NodeOS, please follow these steps:

1. **Do not open a public issue.** This is a security project, and we prioritize responsible disclosure.
2. Send a detailed report to the security contact (as defined in your organizational manifest).
3. Include a Proof of Concept (PoC) if possible.

### What We Care About:
- **Bypassing the Genesis Anchor**: Methods to trick the system into accepting a fake hardware ID.
- **Manifest Forgery**: Any way to modify `manifest.json` without invalidating the Ed25519 signature.
- **Lockdown Evasion**: Techniques to maintain network connectivity after a Lockdown event has been triggered.
- **Identity Leakage**: Potential memory leaks or side-channel attacks that could expose the private node key.

## Response Process

- **Initial Response**: 24-48 hours.
- **Vulnerability Analysis**: 5-7 business days.
- **Patch Release**: Dependent on the severity. Critical vulnerabilities will trigger an immediate emergency update cycle.

---

**SafeStack Security Engineering - Zero-Trust Division**

