# NodeOS – Installation Guide

This guide describes how to install and run **SafeStack NodeOS** (GENESIS Core).
The system is built in Go and follows a **fail-closed, integrity-first** architecture.

---

## Requirements

- **Go ≥ 1.22** (Strictly required)
- **Linux** (tested on Pi 5 / Kali / Ubuntu)
- **Telegram Bot** (for alerts)
- **Root access** (required for Firewall/Maintenance features)

---

## 1. Preparation

Clone the repository and enter the directory:

```bash
cd node-os
```

Install necessary Go dependencies:

```bash
make setup
```

---

## 2. Configuration (Crucial)

Before running, you must provide the Telegram credentials for the awareness engine to work.

1. Create the config file from the example:
   ```bash
   cp config/nodeos_config.example.json config/nodeos_config.json
   ```

2. Edit `config/nodeos_config.json` and insert your tokens.

3. (Optional) Provide a signed manifest in `config/manifest.json`.

---

## 3. Launching the Node

To build and launch the node for the first time:

```bash
make run
```

---

## 4. Maintenance Commands (Makefile)

- `make logs` - Tail the signed crypto-log file.
- `make status` - Show systemd service status.
- `make rebirth` - Clear identity (forces new ID/Keys/Hash on next run).

---
*SAFE · SILENT · FAIL-SAFE*
