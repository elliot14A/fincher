package scheduler

import (
	"math"
	"strings"

	"github.com/elliot14A/fincher/pkg/domain/models"
)

// BaseFailureRate defines the compile-time base failure probability per QC component type.
// These in-code constants are NOT environment-configurable.
var BaseFailureRate = map[models.ComponentType]float64{
	models.ComponentAudio:    0.40,
	models.ComponentSubtitle: 0.25,
	models.ComponentVideo:    0.15,
	models.ComponentMetadata: 0.10,
}

// DefaultBaseFailureRate is the fallback failure rate for unmapped component types.
const DefaultBaseFailureRate = 0.30

// QCOutcome represents the result of a scheduled QC repair inspection.
type QCOutcome string

const (
	QCOutcomePass QCOutcome = "PASS"
	QCOutcomeFail QCOutcome = "FAIL"
)

// DecideOutcome determines whether a completed repair task passes QC inspection.
// If force is "PASSED" or "FAILED", it overrides RNG unconditionally.
// Otherwise, it computes P(fail) = clamp(BaseFailureRate[component] * runVariance, 0.0, 1.0)
// and draws from the scheduler's seeded RNG under mutex.
func (s *Scheduler) DecideOutcome(force string, component models.ComponentType) QCOutcome {
	if strings.EqualFold(force, "PASSED") {
		return QCOutcomePass
	}
	if strings.EqualFold(force, "FAILED") {
		return QCOutcomeFail
	}

	if s == nil {
		return QCOutcomePass
	}

	base, ok := BaseFailureRate[component]
	if !ok {
		base = DefaultBaseFailureRate
	}

	// Clamp effective probability: P(fail) in [0.0, 1.0]
	p := math.Min(1.0, math.Max(0.0, base*s.runVariance))

	s.rngMu.Lock()
	roll := s.rng.Float64()
	s.rngMu.Unlock()

	if roll < p {
		return QCOutcomeFail
	}
	return QCOutcomePass
}

// DefectEventTypeFor maps a defective media component to its domain-scoped CloudEvent type and severity.
func DefectEventTypeFor(component models.ComponentType) (eventType string, severity models.EventSeverity) {
	switch component {
	case models.ComponentAudio:
		return models.TypeAudioSyncDriftDetected, models.SeverityWarn
	case models.ComponentSubtitle:
		return models.TypePackageInvalidated, models.SeverityWarn
	case models.ComponentVideo:
		return models.TypePackageInvalidated, models.SeverityWarn
	case models.ComponentMetadata:
		return models.TypePackageInvalidated, models.SeverityWarn
	default:
		return models.TypePackageInvalidated, models.SeverityWarn
	}
}
