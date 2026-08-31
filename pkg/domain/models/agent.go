package models

// UnmeasuredHistoricalAccuracy indicates that no historical QC records exist for a vendor component.
const UnmeasuredHistoricalAccuracy = -1.0

// AnalyticsSummary captures past vendor performance and defect history from ClickHouse.
type AnalyticsSummary struct {
	VendorHistoricalAccuracy float64  `json:"vendor_historical_accuracy"`
	SimilarDefectOccurrences int      `json:"similar_defect_occurrences"`
	PriorIncidentsForVendor  int      `json:"prior_incidents_for_vendor"`
	RelevantHistoricalLogs   []string `json:"relevant_historical_logs"`
}

// DeliveryImpact captures downstream dependency blast radius from Turso SQLite.
type DeliveryImpact struct {
	RootPackageID      string   `json:"root_package_id"`
	AffectedPackages   []string `json:"affected_packages"`
	AffectedDeliveries []string `json:"affected_deliveries"`
	AffectedMarkets    []string `json:"affected_markets"`
	HoursUntilPremiere float64  `json:"hours_until_premiere"`
	IsPremiereUrgent   bool     `json:"is_premiere_urgent"`
}

// VendorCandidate captures rate, turnaround, and accuracy for vendor selection.
type VendorCandidate struct {
	VendorID           string   `json:"vendor_id"`
	VendorName         string   `json:"vendor_name"`
	Components         []string `json:"components"`
	Markets            []string `json:"markets"`
	HourlyRateUSD      float64  `json:"hourly_rate_usd"`
	TurnaroundHours    int      `json:"turnaround_hours"`
	HistoricalAccuracy float64  `json:"historical_accuracy"`
}

// AllocationRequirement is one component/market needing a vendor in a title-wide allocation plan.
type AllocationRequirement struct {
	Component string `json:"component"` // VIDEO | AUDIO | SUBTITLE | METADATA
	Market    string `json:"market"`    // e.g. "de-DE", "te-IN", or "" for global VIDEO
	Language  string `json:"language"`  // e.g. "de-DE"
}

// RequirementAssignment is the vendor selection for a specific requirement.
type RequirementAssignment struct {
	Component        string  `json:"component"`
	Market           string  `json:"market"`
	Language         string  `json:"language"`
	WinnerVendorID   string  `json:"winner_vendor_id"`
	WinnerVendorName string  `json:"winner_vendor_name"`
	HourlyRateUSD    float64 `json:"hourly_rate_usd"`
	TurnaroundHours  int     `json:"turnaround_hours"`
	Rationale        string  `json:"rationale"`
}

// AllocationPlan represents the holistic staffing decision across all requirements of a title.
type AllocationPlan struct {
	Assignments    []RequirementAssignment `json:"assignments"`
	OverallSummary string                  `json:"overall_summary"`
}
