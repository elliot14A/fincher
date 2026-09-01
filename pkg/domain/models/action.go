package models

// ActionType defines the category of an operational or communication action.
type ActionType string

const (
	ActionHoldDelivery    ActionType = "HOLD_DELIVERY"
	ActionReleaseDelivery ActionType = "RELEASE_DELIVERY"
	ActionReassignVendor  ActionType = "REASSIGN_VENDOR"
	ActionHoldTitle       ActionType = "HOLD_TITLE"

	ActionEmailVendor        ActionType = "EMAIL_VENDOR"
	ActionNotifyStakeholders ActionType = "NOTIFY_STAKEHOLDERS"
	ActionPostSocialUpdate   ActionType = "POST_SOCIAL_UPDATE"
)

// Action represents an individual task proposed by an agent optimizer.
type Action struct {
	Type     ActionType     `json:"type" validate:"required"`
	TargetID string         `json:"target_id" validate:"required"`
	Reason   string         `json:"reason" validate:"required"`
	Payload  map[string]any `json:"payload,omitempty"`
}

// ActionPlan encapsulates a sequence of proposed actions for a title incident.
type ActionPlan struct {
	TitleSlug string   `json:"title_slug"`
	Summary   string   `json:"summary"`
	Actions   []Action `json:"actions"`
}

// Validate verifies ActionPlan constraints.
func (p *ActionPlan) Validate() error {
	if p.TitleSlug == "" {
		p.TitleSlug = DefaultTitleAgnosticSentinel
	}
	for i := range p.Actions {
		if err := validate.Struct(&p.Actions[i]); err != nil {
			return err
		}
	}
	return nil
}
