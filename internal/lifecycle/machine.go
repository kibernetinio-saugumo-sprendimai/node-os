package lifecycle

import "errors"

type Machine struct {
	current State
}

func New() *Machine {
	return &Machine{current: Init}
}

func (m *Machine) State() State {
	return m.current
}

func (m *Machine) Transition(next State) error {
	if !allowed(m.current, next) {
		return errors.New("illegal lifecycle transition")
	}
	m.current = next
	return nil
}

func allowed(from, to State) bool {
	switch from {
	case Init:
		return to == IdentityBorn
	case IdentityBorn:
		return to == Aware || to == Lockdown || to == Terminated
	case Aware:
		return to == Degraded || to == Lockdown || to == Terminating
	case Degraded:
		return to == Aware || to == Lockdown || to == Terminating
	case Lockdown:
		return to == Aware || to == Terminating
	case Terminating:
		return to == Terminated
	default:
		return false
	}
}
