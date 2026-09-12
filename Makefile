# SafeStack NodeOS - Makefile
# GENESIS Core (v0.2.0-genesis)

BINARY_NAME=nodeos

all: build

# Build the node binary
build:
	@echo "Building SafeStack NodeOS (Optimized)..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) nodeos.go

# Run the node
run: build:
	@echo "Launching NodeOS..."
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	@echo "Cleaning binaries..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# Full Clean: Wipe everything including logs and identity
full-clean: clean
	@echo "Performing DEEP CLEAN..."
	rm -f node_id.txt node_key.txt node_bin.hash genesis_hash.txt
	rm -f config/manifest.json.sig SHA256SUMS SHA256SUMS.sig
	rm -f AUDIT_REPORT.md AUDIT_REPORT.md.sig
	rm -f /var/lib/nodeos/LOCKDOWN /var/lib/nodeos/nodeos.crypt.log

# REBIRTH: Wipe identity and force new generation
rebirth:
	@echo "Wiping identity and binary hash..."
	rm -f node_id.txt node_key.txt node_bin.hash genesis_hash.txt
	rm -f /var/lib/nodeos/LOCKDOWN
	@echo "Ready for REBIRTH on next run."

# Setup dependencies
setup:
	@echo "Preparing environment..."
	@echo "NOTE: Ensure 'smartmontools' is installed for NVMe Health (sudo apt install smartmontools)"
	go mod tidy
	mkdir -p config docs systemd

# Run internal tests
test:
	@echo "Running package tests..."
	go test -v ./internal/...

# Verify project integrity. Any mismatch must return a non-zero exit code.
verify:
	@echo "Checking project checksums..."
	@sha256sum -c SHA256SUMS

# Sign files using a path to a mode-0600 Ed25519 private key.
# Usage: make sign-manifest KEY_FILE=/secure/path/manifest.key
sign-manifest:
	@test -n "$(KEY_FILE)" || (echo "KEY_FILE is required" >&2; exit 2)
	@go run scripts/sign_manifest.go "$(KEY_FILE)" config/manifest.json

# Usage: make sign-audit KEY_FILE=/secure/path/audit.key
sign-audit:
	@test -n "$(KEY_FILE)" || (echo "KEY_FILE is required" >&2; exit 2)
	@go run scripts/sign_manifest.go "$(KEY_FILE)" AUDIT_REPORT.md SHA256SUMS

# Help
help:
	@echo "SafeStack NodeOS - Command Registry"
	@echo ""
	@echo "Usage:"
	@echo "  make build         - Compile the node executable (optimized)"
	@echo "  make run           - Compile and launch the node"
	@echo "  make setup         - Prepare environment and dependencies"
	@echo "  make test          - Run unit tests"
	@echo "  make verify        - Verify SHA256SUMS and fail on mismatch"
	@echo "  make sign-manifest KEY_FILE=/secure/path/key"
	@echo "  make sign-audit    KEY_FILE=/secure/path/key"
	@echo "  make rebirth       - Reset identity (preserves logs)"
	@echo "  make full-clean    - Wipe identity, logs and audit artifacts"
	@echo "  make clean         - Remove binary artifacts"
