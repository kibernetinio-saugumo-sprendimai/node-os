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
- **Tapatybės ciklas:** Tikrina `Rebirth`, `Wipe` ir `Init` funkcijas.
- **Failų generavimas:** Užtikrina, kad raktai ir ID būtų teisingai išsaugoti diske su tinkamomis teisėmis.

### 2. Configuration & Manifest (`internal/config`)
- **Parašų validavimas:** Tikrina, ar sistema teisingai atpažįsta galiojančius ir suklastotus manifesto parašus.
- **Root of Trust:** Užtikrina, kad kietai įrašytas viešasis raktas veikia kaip pagrindinis pasitikėjimo šaltinis.

## Security Auditing

Prieš kiekvieną *release*, rekomenduojama rankiniu būdu patikrinti manifesto vientisumą:
1. Pakeiskite reikšmę `config/manifest.json`.
2. Paleiskite `make test`.
3. Įsitikinkite, kad parašų tikrinimas nepraeina be naujo pasirašymo.
