package models

import "time"

// Vendor represents a post-production studio delivering media assets.
type Vendor struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Specialty string    `json:"specialty" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate verifies vendor attributes.
func (v *Vendor) Validate() error {
	return validate.Struct(v)
}

// UpdateVendorInput represents partial update fields for a Vendor.
type UpdateVendorInput struct {
	Name      *string `json:"name,omitempty" validate:"omitempty,min=1"`
	Specialty *string `json:"specialty,omitempty" validate:"omitempty,min=1"`
}

// Validate verifies vendor update constraints.
func (u *UpdateVendorInput) Validate() error {
	return validate.Struct(u)
}
