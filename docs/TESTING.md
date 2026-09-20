# NodeOS Testing Protocols

To maintain high security standards, NodeOS includes internal unit tests for critical cryptographic and identity functions.

## Running Tests

Tests are managed via the standard Go test tool. You can run all package tests using the Makefile:

```bash
make test
```

Or directly using `go`:

```bash
go test ./internal/...
```

## Test Coverage

### 1. Identity Layer (`internal/identity`)
- **Identity lifecycle:** Tests the `Rebirth`, `Wipe` and `Init` functions.
- **File generation:** Ensures that keys and IDs are saved correctly on disk with the appropriate permissions.

### 2. Configuration & Manifest (`internal/config`)
- **Signature validation:** Tests whether the system correctly recognizes valid and forged manifest signatures.
- **Root of Trust:** Ensures that the hard-coded public key acts as the primary trust anchor.

## Security Auditing

Before each *release*, manually verify manifest integrity:
1. Change a value in `config/manifest.json`.
2. Run `make test`.
3. Confirm that signature verification fails until the manifest is signed again.
