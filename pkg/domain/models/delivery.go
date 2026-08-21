package models

import "time"

// DeliveryStatus represents the shipping state of a territory delivery.
type DeliveryStatus string

const (
	DeliveryStatusPending     DeliveryStatus = "PENDING"
	DeliveryStatusReadyToShip DeliveryStatus = "READY_TO_SHIP"
	DeliveryStatusHold        DeliveryStatus = "HOLD"
	DeliveryStatusShipped     DeliveryStatus = "SHIPPED"
)

// Delivery represents a territory release target.
type Delivery struct {
	Base
	TitleID    string         `json:"title_id" validate:"required"`
	Country    string         `json:"country" validate:"required"`
	Status     DeliveryStatus `json:"status" validate:"required,oneof=PENDING READY_TO_SHIP HOLD SHIPPED"`
	TargetDate time.Time      `json:"target_date" validate:"required"`
}

// Validate verifies delivery constraints.
func (d *Delivery) Validate() error {
	if err := d.ValidateMetadata(); err != nil {
		return err
	}
	return validate.Struct(d)
}

// UpdateDeliveryInput represents partial update fields for a Delivery.
type UpdateDeliveryInput struct {
	Country    *string         `json:"country,omitempty" validate:"omitempty,min=1"`
	Status     *DeliveryStatus `json:"status,omitempty" validate:"omitempty,oneof=PENDING READY_TO_SHIP HOLD SHIPPED"`
	TargetDate *time.Time      `json:"target_date,omitempty"`
	Metadata   map[string]any  `json:"metadata,omitempty"`
}

// Validate verifies partial delivery update constraints.
func (u *UpdateDeliveryInput) Validate() error {
	if err := ValidateMetadataMap(u.Metadata); err != nil {
		return err
	}
	return validate.Struct(u)
}
