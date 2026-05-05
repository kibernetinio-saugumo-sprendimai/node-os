package policy

const MaxAutonomy = 5

func CheckAutonomy(level int) Signal {
	if level < 0 || level > MaxAutonomy {
		return Signal{
			Severity:   SeverityCritical,
			Reason:     ReasonAutonomyOutOfRange,
			Confidence: 1.0,
		}
	}

	return Signal{
		Severity:   SeverityOK,
		Confidence: 1.0,
	}
}
