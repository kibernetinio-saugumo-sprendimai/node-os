# Deployment Guide: Desktop to Raspberry Pi 5
> **Securely transferring SafeStack NodeOS to your edge hardware.**

This guide explains how to deploy the NodeOS codebase from your development machine (Windows/PowerShell) to your Raspberry Pi 5 (Argon NEO 5 NVMe) using Secure Copy (SCP).

---

## 🏗️ Deployment Strategy
Since SafeStack NodeOS is a security-focused project, we avoid unencrypted transfers. We use **SCP over SSH** to maintain end-to-end encryption during the deployment process.

### Prerequisites
- **Raspberry Pi 5** reachable over the network.
- **SSH enabled** on the Pi.
- **PowerShell 7+** or Windows PowerShell with `OpenSSH` installed.

---

## 🚀 PowerShell Deployment Script

Run the following script from the root of the project directory on your Windows machine:

```powershell
# --- CONFIGURATION ---
$PI_IP   = "192.168.1.XXX"  # Change to your Pi's actual IP
$PI_USER = "pi"             # Change to your username
$PORT    = 22               # Default is 22 (Change to 22009 if already hardened)

# --- EXECUTION ---
Write-Host "📡 Connecting to Pi at $PI_IP..." -ForegroundColor Cyan

# 1. Create target directory
ssh -p $PORT ${PI_USER}@${PI_IP} "mkdir -p ~/node-os"

# 2. Transfer files (excluding local sensitive state)
# Note: .gitignore will be ignored by scp, so we transfer everything.
# Rebirth on the Pi will generate new local identity files.
scp -P $PORT -r ./* ${PI_USER}@${PI_IP}:~/node-os/

Write-Host "✅ Deployment Complete!" -ForegroundColor Green
Write-Host "Next Step: SSH into the Pi and run 'make setup'" -ForegroundColor Yellow
```

---

## 🛠️ Post-Transfer Initialization

Once the files are on the Pi, connect via SSH and perform the initial setup:

```bash
# 1. Connect to the Pi
ssh -p 22 pi@192.168.1.XXX

# 2. Navigate to the directory
cd ~/node-os

# 3. Setup dependencies and build
make setup
make build

# 4. Initial Launch (Identity Generation)
make run
```

---

## ⚓ Identity Persistence
Once you run `make run` for the first time on the Pi 5, the following files will be created locally on the device:
- `node_id.txt`
- `node_key.txt`
- `genesis_hash.txt` (Locked as immutable)

**Do not transfer these files back to your PC.** They are bound to your Argon NEO 5 hardware. Future deployments using this guide will only update the code, preserving your node's unique identity.

**SafeStack Security Engineering - Deployment Operations**
