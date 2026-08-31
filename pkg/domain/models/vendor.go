package models

import "fmt"

// Vendor represents a post-production localization or QC facility.
type Vendor struct {
	Base
	Name            string   `json:"name" validate:"required"`
	Components      []string `json:"components" validate:"required,min=1"`
	Markets         []string `json:"markets"`
	HourlyRateUSD   float64  `json:"hourly_rate_usd" validate:"gte=0"`
	TurnaroundHours int      `json:"turnaround_hours" validate:"gt=0"`
}

// Validate verifies vendor constraints.
func (v *Vendor) Validate() error {
	if err := v.ValidateMetadata(); err != nil {
		return err
	}
	if v.TurnaroundHours == 0 {
		v.TurnaroundHours = 24
	}
	if len(v.Components) == 0 {
		return fmt.Errorf("vendor must support at least one component")
	}
	for _, comp := range v.Components {
		switch comp {
		case string(ComponentVideo), string(ComponentAudio), string(ComponentSubtitle), string(ComponentMetadata):
			// valid
		default:
			return fmt.Errorf("invalid component: %s", comp)
		}
	}
	return validate.Struct(v)
}

// UpdateVendorInput represents partial update fields for a Vendor.
type UpdateVendorInput struct {
	Name            *string        `json:"name,omitempty" validate:"omitempty,min=1"`
	Components      *[]string      `json:"components,omitempty" validate:"omitempty,min=1"`
	Markets         *[]string      `json:"markets,omitempty"`
	HourlyRateUSD   *float64       `json:"hourly_rate_usd,omitempty" validate:"omitempty,gte=0"`
	TurnaroundHours *int           `json:"turnaround_hours,omitempty" validate:"omitempty,gt=0"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

// Validate verifies partial vendor update constraints.
func (u *UpdateVendorInput) Validate() error {
	if err := ValidateMetadataMap(u.Metadata); err != nil {
		return err
	}
	if u.Components != nil {
		if len(*u.Components) == 0 {
			return fmt.Errorf("components cannot be empty")
		}
		for _, comp := range *u.Components {
			switch comp {
			case string(ComponentVideo), string(ComponentAudio), string(ComponentSubtitle), string(ComponentMetadata):
				// valid
			default:
				return fmt.Errorf("invalid component: %s", comp)
			}
		}
	}
	return validate.Struct(u)
}
