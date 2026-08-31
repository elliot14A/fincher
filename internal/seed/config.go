package seed

import "github.com/elliot14A/fincher/internal/seed/types"

// SeedConfig aliases the subsystem configuration struct.
type SeedConfig = types.SeedConfig

// DefaultConfig returns a SeedConfig initialized with default values.
func DefaultConfig() *SeedConfig {
	return types.DefaultConfig()
}
