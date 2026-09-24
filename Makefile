# SafeStack NodeOS - Makefile
# GENESIS Core (v0.2.0-genesis)

BINARY_NAME=nodeos

all: build

# 🔨 Build the node binary
build:
	@echo "🔨 Building SafeStack NodeOS (Optimized)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) nodeos.go

# 🚀 Run the node
run: build
	@echo "🚀 Launching NodeOS..."
	./$(BINARY_NAME)

# 🧹 Clean build artifacts
clean:
	@echo "🧹 Cleaning binaries..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# 🧼 Full Clean: Wipe everything including logs and identity
full-clean: clean
	@echo "🧼 Performing DEEP CLEAN..."
	rm -f node_id.txt node_key.txt node_bin.hash genesis_hash.txt
	rm -f /var/lib/nodeos/LOCKDOWN /var/lib/nodeos/nodeos.crypt.log

# 🧬 REBIRTH: Wipe identity and force new generation
rebirth:
	@echo "⚠️  Wiping identity and binary hash..."
	rm -f node_id.txt node_key.txt node_bin.hash genesis_hash.txt
	rm -f /var/lib/nodeos/LOCKDOWN
	@echo "✨ Ready for REBIRTH on next run."

# 🛠️ Setup dependencies
setup:
	@echo "🛠️ Preparing environment..."
	@echo "📦 NOTE: Ensure 'smartmontools' is installed for NVMe Health (sudo apt install smartmontools)"
	go mod tidy
	mkdir -p config docs systemd

# 🧪 Run internal tests
test:
	@echo "🧪 Running package tests..."
	go test -v ./internal/...

# 🛡️ Verify project integrity
verify:
	@echo "🛡️ Checking project signatures..."
	@go run scripts/verify/main.go
	@echo "🛡️ Verifying SHA256 file hashes..."
	@(command -v sha256sum >/dev/null 2>&1 && sha256sum -c SHA256SUMS) || shasum -a 256 -c SHA256SUMS || echo "❌ INTEGRITY CHECK FAILED"

# 🖋️ Sign the manifest (Usage: make sign-manifest KEY=<priv_key>)
sign-manifest:
	@echo "🖋️ Signing manifest..."
	go run scripts/sign_manifest.go $(KEY) config/manifest.json

# 📜 Sign the Audit Trail (Usage: make sign-audit KEY=<priv_key>)
sign-audit:
	@echo "📜 Signing Audit Trail and Checksums..."
	@go run -e "import ('crypto/ed25519'; 'encoding/hex'; 'os'); func main() { p,_ := hex.DecodeString(\"$(KEY)\"); priv := ed25519.PrivateKey(p); if len(p)==32 { priv = ed25519.NewKeyFromSeed(p) }; for _, f := range []string{\"AUDIT_REPORT.md\", \"SHA256SUMS\"} { d,_ := os.ReadFile(f); s := ed25519.Sign(priv, d); os.WriteFile(f+\".sig\", []byte(hex.EncodeToString(s)), 0644) } }"

# 📖 Help
help:
	@echo "SafeStack NodeOS - Command Registry"
	@echo ""
	@echo "Usage:"
	@echo "  make build         - Compile the node executable (optimized)"
	@echo "  make run           - Compile and launch the node"
	@echo "  make setup         - Prepare environment and dependencies"
	@echo "  make test          - Run unit tests"
	@echo "  make verify        - Verify project integrity (SHA256SUMS & signatures)"
	@echo "  make sign-manifest - Sign manifest.json (KEY=<hex>)"
	@echo "  make sign-audit    - Sign Audit Report & Checksums (KEY=<hex>)"
	@echo "  make rebirth       - Reset identity (preserves logs)"
	@echo "  make full-clean    - Wipe EVERYTHING (identity, logs, audit)"
	@echo "  make clean         - Remove binary artifacts"
