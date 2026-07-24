package health

import "strings"

// HealthState is the reachability lifecycle for a catalog model.
type HealthState string

const (
	// StateHealthy means the latest probe passed and error rate is below threshold.
	StateHealthy HealthState = "healthy"
	// StateUnknown means the model has never been probed (or was reset).
	StateUnknown HealthState = "unknown"
	// StateDegraded means the probe passed but live error rate exceeds the unhealthy threshold.
	StateDegraded HealthState = "degraded"
	// StateUnhealthy means probes failed and the model is not advertised.
	StateUnhealthy HealthState = "unhealthy"
	// StateRecovering means probes failed and exponential backoff is scheduling retries.
	StateRecovering HealthState = "recovering"
)

// IsAdvertisable reports whether a state may appear in /v1/models.
// Degraded stays visible (probe passed; elevated live error rate is advisory).
// Recovering and Unhealthy are hidden when hide_unreachable is enabled.
func (s HealthState) IsAdvertisable(unknownAsReachable bool) bool {
	switch s {
	case StateHealthy, StateDegraded:
		return true
	case StateUnknown, "":
		return unknownAsReachable
	default:
		return false
	}
}

// DeriveState maps legacy Reachable/ConsecutiveFails into a HealthState when
// State was not persisted (older DB rows).
func DeriveState(reachable bool, consecutiveFails int, known bool) HealthState {
	if !known {
		return StateUnknown
	}
	if reachable {
		return StateHealthy
	}
	if consecutiveFails > 0 {
		return StateRecovering
	}
	return StateUnhealthy
}

// NormalizeState returns a known HealthState, defaulting empty to unknown.
func NormalizeState(s HealthState) HealthState {
	switch HealthState(strings.ToLower(string(s))) {
	case StateHealthy:
		return StateHealthy
	case StateDegraded:
		return StateDegraded
	case StateUnhealthy:
		return StateUnhealthy
	case StateRecovering:
		return StateRecovering
	case StateUnknown:
		return StateUnknown
	default:
		return StateUnknown
	}
}
