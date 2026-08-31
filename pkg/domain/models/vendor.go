package models

// Vendor represents a post-production localization or QC facility.
type Vendor struct {
	Base
	Name            string  `json:"name" validate:"required"`
	Specialty       string  `json:"specialty" validate:"required"`
	HourlyRateUSD   float64 `json:"hourly_rate_usd" validate:"gte=0"`
	TurnaroundHours int     `json:"turnaround_hours" validate:"gt=0"`
}

// Validate verifies vendor constraints.
func (v *Vendor) Validate() error {
	if err := v.ValidateMetadata(); err != nil {
		return err
	}
	if v.TurnaroundHours == 0 {
		v.TurnaroundHours = 24
	}
	return validate.Struct(v)
}

// UpdateVendorInput represents partial update fields for a Vendor.
type UpdateVendorInput struct {
	Name            *string        `json:"name,omitempty" validate:"omitempty,min=1"`
	Specialty       *string        `json:"specialty,omitempty" validate:"omitempty,min=1"`
	HourlyRateUSD   *float64       `json:"hourly_rate_usd,omitempty" validate:"omitempty,gte=0"`
	TurnaroundHours *int           `json:"turnaround_hours,omitempty" validate:"omitempty,gt=0"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

// Validate verifies partial vendor update constraints.
func (u *UpdateVendorInput) Validate() error {
	if err := ValidateMetadataMap(u.Metadata); err != nil {
		return err
	}
	return validate.Struct(u)
}
