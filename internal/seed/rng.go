package seed

import "github.com/elliot14A/fincher/internal/seed/types"

// RNG aliases the deterministic RNG wrapper.
type RNG = types.RNG

// NewRNG constructs an RNG initialized with the provided deterministic seed.
func NewRNG(seed int64) *RNG {
	return types.NewRNG(seed)
}

// Choice returns a random element from a non-empty slice.
func Choice[T any](r *RNG, items []T) T {
	return types.Choice(r, items)
}
