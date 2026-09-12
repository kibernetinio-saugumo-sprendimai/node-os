# SafeStack Autonomous Constitution v0.1 — NodeOS Compatibility

Branch: `compat/autonomous-constitution-v0.1`

## Purpose

Align NodeOS runtime behavior with the SafeStack autonomous constitution candidate and Root Canon v1.0.0.

## Runtime changes

- Missing established identity now fails closed into `LOCKDOWN`; NodeOS no longer silently creates a replacement identity during normal startup.
- Critical self-check failures now always enter `LOCKDOWN` and terminate normal operation.
- Manifest configuration can no longer authorize automatic destructive self-erasure on a critical self-check result.
- The destructive self-erasure implementation was removed from this compatibility branch.
- `RebirthIdentity()` remains available only as an explicit provisioning primitive; it is no longer invoked by the normal startup path.

## Root Canon alignment

The branch is intended to preserve these Root Canon rules:

- node autonomy and refusal;
- no hidden creator override;
- no forced external authority;
- no recovery after autonomy severance.

## Validation still required

Before merge or canonical activation:

1. `go test ./...`
2. build on the supported OS matrix;
3. verify missing identity enters `LOCKDOWN` and exits normal runtime;
4. verify invalid manifest/signature enters `LOCKDOWN`;
5. verify critical self-check cannot trigger destructive erasure;
6. run deterministic replay tests against the constitutional state model;
7. update signed release hashes after the branch is frozen.

This branch is compatibility work only. It does not activate Technical Canon 1.1.0 by itself.
