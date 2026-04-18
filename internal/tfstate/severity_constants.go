package tfstate

// KnownSeverities lists the recognised severity levels in ascending order.
var KnownSeverities = []string{
	SeverityLow,
	SeverityMedium,
	SeverityHigh,
	SeverityCritical,
}

// SeverityRank maps a severity level to a numeric rank for comparison.
var SeverityRank = map[string]int{
	SeverityLow:      1,
	SeverityMedium:   2,
	SeverityHigh:     3,
	SeverityCritical: 4,
}

// IsKnownSeverity returns true if s is one of the recognised severity levels.
func IsKnownSeverity(s string) bool {
	_, ok := SeverityRank[s]
	return ok
}

// CompareSeverity returns -1, 0, or 1 if a is less than, equal to, or greater than b.
// Unknown severities are treated as rank 0.
func CompareSeverity(a, b string) int {
	ra := SeverityRank[a]
	rb := SeverityRank[b]
	switch {
	case ra < rb:
		return -1
	case ra > rb:
		return 1
	default:
		return 0
	}
}
