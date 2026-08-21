package models

import (
	"strings"
	"time"
)

// ComponentType denotes the variety of media asset.
type ComponentType string

const (
	ComponentVideo    ComponentType = "VIDEO"
	ComponentAudio    ComponentType = "AUDIO"
	ComponentSubtitle ComponentType = "SUBTITLE"
	ComponentMetadata ComponentType = "METADATA"
)

// PackageStatus indicates component validation and release state.
type PackageStatus string

const (
	PackageStatusPending     PackageStatus = "PENDING"
	PackageStatusValid       PackageStatus = "VALID"
	PackageStatusInvalidated PackageStatus = "INVALIDATED"
	PackageStatusReQCPending PackageStatus = "RE_QC_PENDING"
)

// Package represents a deliverable media component asset.
type Package struct {
	ID                       string        `json:"id" validate:"required"`
	TitleID                  string        `json:"title_id" validate:"required"`
	Component                ComponentType `json:"component" validate:"required,oneof=VIDEO AUDIO SUBTITLE METADATA"`
	Language                 string        `json:"language" validate:"required"`
	Version                  string        `json:"version" validate:"required"`
	VendorID                 string        `json:"vendor_id" validate:"required"`
	DerivedFromMasterVersion string        `json:"derived_from_master_version" validate:"required"`
	RedeliveryCount          int           `json:"redelivery_count" validate:"gte=0"`
	Status                   PackageStatus `json:"status" validate:"required,oneof=PENDING VALID INVALIDATED RE_QC_PENDING"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
}

// Validate verifies package constraints.
func (p *Package) Validate() error {
	return validate.Struct(p)
}

// IsStaleAgainst checks if package's master version differs from active master version.
func (p *Package) IsStaleAgainst(activeMasterVersion string) bool {
	return strings.TrimSpace(p.DerivedFromMasterVersion) != strings.TrimSpace(activeMasterVersion)
}

// UpdatePackageInput represents partial update fields for a Package.
type UpdatePackageInput struct {
	Component                *ComponentType `json:"component,omitempty" validate:"omitempty,oneof=VIDEO AUDIO SUBTITLE METADATA"`
	Language                 *string        `json:"language,omitempty" validate:"omitempty,min=1"`
	Version                  *string        `json:"version,omitempty" validate:"omitempty,min=1"`
	VendorID                 *string        `json:"vendor_id,omitempty" validate:"omitempty,min=1"`
	DerivedFromMasterVersion *string        `json:"derived_from_master_version,omitempty" validate:"omitempty,min=1"`
	RedeliveryCount          *int           `json:"redelivery_count,omitempty" validate:"omitempty,gte=0"`
	Status                   *PackageStatus `json:"status,omitempty" validate:"omitempty,oneof=PENDING VALID INVALIDATED RE_QC_PENDING"`
}

// Validate verifies partial package update constraints.
func (u *UpdatePackageInput) Validate() error {
	return validate.Struct(u)
}
