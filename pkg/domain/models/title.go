package models

import "time"

// TitleStatus represents the release readiness of a title.
type TitleStatus string

const (
	StatusDraft      TitleStatus = "DRAFT"
	StatusOnTrack    TitleStatus = "ON_TRACK"
	StatusAtRisk     TitleStatus = "AT_RISK"
	StatusHold       TitleStatus = "HOLD"
	StatusProcessing TitleStatus = "PROCESSING"
	StatusShipped    TitleStatus = "SHIPPED"
	StatusOverdue    TitleStatus = "OVERDUE"
)

// TitleType denotes the category of media release.
type TitleType string

const (
	TitleTypeFeature TitleType = "FEATURE"
	TitleTypeSeries  TitleType = "SERIES"
	TitleTypeSpecial TitleType = "SPECIAL"
)

// Title represents a media release entity.
type Title struct {
	Base
	Name                 string      `json:"name" validate:"required"`
	Slug                 string      `json:"slug" validate:"required"`
	Type                 TitleType   `json:"type" validate:"required,oneof=FEATURE SERIES SPECIAL"`
	PremiereDate         time.Time   `json:"premiere_date" validate:"required"`
	Territories          int         `json:"territories" validate:"required,gte=1"`
	CurrentMasterVersion string      `json:"current_master_version" validate:"required"`
	OverallStatus        TitleStatus `json:"overall_status" validate:"required,oneof=DRAFT ON_TRACK AT_RISK HOLD PROCESSING SHIPPED OVERDUE"`
}

// Validate verifies title attributes.
func (t *Title) Validate() error {
	if err := t.ValidateMetadata(); err != nil {
		return err
	}
	return validate.Struct(t)
}

// UpdateTitleInput represents partial update attributes for a Title.
type UpdateTitleInput struct {
	Name                 *string        `json:"name,omitempty" validate:"omitempty,min=1"`
	Slug                 *string        `json:"slug,omitempty" validate:"omitempty,min=1"`
	Type                 *TitleType     `json:"type,omitempty" validate:"omitempty,oneof=FEATURE SERIES SPECIAL"`
	PremiereDate         *time.Time     `json:"premiere_date,omitempty"`
	Territories          *int           `json:"territories,omitempty" validate:"omitempty,gte=1"`
	CurrentMasterVersion *string        `json:"current_master_version,omitempty" validate:"omitempty,min=1"`
	OverallStatus        *TitleStatus   `json:"overall_status,omitempty" validate:"omitempty,oneof=DRAFT ON_TRACK AT_RISK HOLD PROCESSING SHIPPED OVERDUE"`
	Metadata             map[string]any `json:"metadata,omitempty"`
}

// Validate verifies partial title update constraints.
func (u *UpdateTitleInput) Validate() error {
	if err := ValidateMetadataMap(u.Metadata); err != nil {
		return err
	}
	return validate.Struct(u)
}

// HoursUntilPremiere calculates remaining duration in hours from the given reference time.
func (t *Title) HoursUntilPremiere(from time.Time) float64 {
	return t.PremiereDate.Sub(from).Hours()
}

// IsImminentLaunch checks if premiere date is within threshold hours from reference time.
func (t *Title) IsImminentLaunch(from time.Time, thresholdHours float64) bool {
	hours := t.HoursUntilPremiere(from)
	return hours >= 0 && hours <= thresholdHours
}
