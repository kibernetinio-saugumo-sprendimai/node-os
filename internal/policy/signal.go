package policy

type Severity int

const (
    SeverityOK Severity = iota
    SeverityWarn
    SeverityCritical
)

func (s Severity) String() string {
    switch s {
    case SeverityOK:
        return "OK"
    case SeverityWarn:
        return "WARN"
    case SeverityCritical:
        return "CRITICAL"
    default:
        return "UNKNOWN"
    }
}

type Reason string

const (
    ReasonAutonomyOutOfRange Reason = "AUTONOMY_OUT_OF_RANGE"
    ReasonBinaryTampered     Reason = "BINARY_TAMPERED"
    ReasonBinaryReadFailed   Reason = "BINARY_READ_FAILED"
    ReasonGenesisMismatch    Reason = "GENESIS_ANCHOR_MISMATCH"
    ReasonManifestDrift      Reason = "MANIFEST_DRIFT"
    ReasonDebuggerDetected   Reason = "DEBUGGER_DETECTED"
    ReasonInsecurePermissions Reason = "INSECURE_FILE_PERMISSIONS"
    ReasonMissingRoot        Reason = "MISSING_ROOT_AUTHORITY"
    ReasonAttackSurfaceHigh  Reason = "ATTACK_SURFACE_HIGH"
)

type Signal struct {
    Severity   Severity
    Reason     Reason
    Confidence float64
}
