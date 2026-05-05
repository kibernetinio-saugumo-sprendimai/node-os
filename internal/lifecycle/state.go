package lifecycle

type State string

const (
	Init         State = "INIT"
	IdentityBorn State = "IDENTITY_BORN"
	Aware        State = "AWARE"
	Degraded     State = "DEGRADED"
	Lockdown     State = "LOCKDOWN"
	Terminating  State = "TERMINATING"
	Terminated   State = "TERMINATED"
)
