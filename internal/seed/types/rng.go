package types

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
)

// RNG provides deterministic pseudo-random number generation and faker utilities.
type RNG struct {
	Rand  *rand.Rand
	Faker *gofakeit.Faker
}

// NewRNG constructs an RNG initialized with the provided deterministic seed.
func NewRNG(seed int64) *RNG {
	r := rand.New(rand.NewSource(seed))
	f := gofakeit.New(uint64(seed))
	return &RNG{
		Rand:  r,
		Faker: f,
	}
}

// Float64 returns a pseudo-random float64 in [0.0, 1.0).
func (r *RNG) Float64() float64 {
	return r.Rand.Float64()
}

// FloatInRange returns a float64 in [min, max).
func (r *RNG) FloatInRange(min, max float64) float64 {
	if min >= max {
		return min
	}
	return min + r.Rand.Float64()*(max-min)
}

// IntInRange returns an int in [min, max] (inclusive).
func (r *RNG) IntInRange(min, max int) int {
	if min >= max {
		return min
	}
	return min + r.Rand.Intn(max-min+1)
}

// NormFloat64 returns a normally distributed float64 with given mean and standard deviation.
func (r *RNG) NormFloat64(mean, stdDev float64) float64 {
	return mean + r.Rand.NormFloat64()*stdDev
}

// Choice returns a random element from a non-empty slice.
func Choice[T any](r *RNG, items []T) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	return items[r.Rand.Intn(len(items))]
}
