package models

// RepairProjection captures in-flight or estimated repair work for a package.
type RepairProjection struct {
	PackageID       string  `json:"package_id"`
	Component       string  `json:"component"`
	Market          string  `json:"market"`
	VendorID        string  `json:"vendor_id"`
	VendorName      string  `json:"vendor_name"`
	TurnaroundHours float64 `json:"turnaround_hours"`
	Status          string  `json:"status"`
}

// TitleProjection provides pre-computed timeline, critical path, buffer hours, and risk band for a title.
type TitleProjection struct {
	TitleSlug              string             `json:"title_slug"`
	TitleName              string             `json:"title_name"`
	PremiereDate           string             `json:"premiere_date"`
	HoursUntilPremiere     float64            `json:"hours_until_premiere"`
	CriticalRemainingHours float64            `json:"critical_remaining_hours"`
	BufferHours            float64            `json:"buffer_hours"`
	RiskBand               string             `json:"risk_band"` // BREACH | TIGHT | WATCH | SAFE
	IsBreached             bool               `json:"is_breached"`
	IsUrgent               bool               `json:"is_urgent"`
	Repairs                []RepairProjection `json:"repairs"`
}
