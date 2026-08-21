package models

// Vendor represents a post-production localization or QC facility.
type Vendor struct {
	Base
	Name      string `json:"name" validate:"required"`
	Specialty string `json:"specialty" validate:"required"`
}

// Validate verifies vendor constraints.
func (v *Vendor) Validate() error {
	if err := v.ValidateMetadata(); err != nil {
		return err
	}
	return validate.Struct(v)
}

// UpdateVendorInput represents partial update fields for a Vendor.
type UpdateVendorInput struct {
	Name      *string        `json:"name,omitempty" validate:"omitempty,min=1"`
	Specialty *string        `json:"specialty,omitempty" validate:"omitempty,min=1"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

// Validate verifies partial vendor update constraints.
func (u *UpdateVendorInput) Validate() error {
	if err := ValidateMetadataMap(u.Metadata); err != nil {
		return err
	}
	return validate.Struct(u)
}
