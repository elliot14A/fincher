package catalog

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/elliot14A/fincher/internal/seed/types"
)

// VendorSpec defines the configuration profile for a curated vendor.
type VendorSpec struct {
	ID              string
	Name            string
	Components      []string
	Markets         []string
	HourlyRateUSD   float64
	TurnaroundHours float64
	TargetAccuracy  float64 // Used for ClickHouse historical event generation
	DriftMeanMS     float64 // Used for synthetic audio sync drift distribution
	DriftStdDevMS   float64
	PosterURL       string
}

func makePosterURL(name, bg, color string) string {
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=%s&color=%s&bold=true&size=128", url.QueryEscape(name), bg, color)
}

// CuratedVendors returns the 8 canonical post-production and localization facilities.
func CuratedVendors() []VendorSpec {
	allFiveMarkets := []string{"en-US", "de-DE", "fr-FR", "hi-IN", "te-IN"}
	westernMarkets := []string{"en-US", "de-DE", "fr-FR"}
	indianMarkets := []string{"hi-IN", "te-IN"}

	return []VendorSpec{
		{
			ID:              "vnd-deluxe",
			Name:            "Deluxe Media",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         allFiveMarkets,
			HourlyRateUSD:   200.0,
			TurnaroundHours: 12.0,
			TargetAccuracy:  0.99,
			DriftMeanMS:     20.0,
			DriftStdDevMS:   10.0,
			PosterURL:       makePosterURL("Deluxe Media", "0284c7", "ffffff"),
		},
		{
			ID:              "vnd-iyuno",
			Name:            "Iyuno SDI",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         westernMarkets,
			HourlyRateUSD:   70.0,
			TurnaroundHours: 36.0,
			TargetAccuracy:  0.93,
			DriftMeanMS:     45.0,
			DriftStdDevMS:   20.0,
			PosterURL:       makePosterURL("Iyuno SDI", "7c3aed", "ffffff"),
		},
		{
			ID:              "vnd-testronic",
			Name:            "Testronic Labs",
			Components:      []string{"AUDIO"},
			Markets:         westernMarkets,
			HourlyRateUSD:   120.0,
			TurnaroundHours: 24.0,
			TargetAccuracy:  0.85,
			DriftMeanMS:     140.0,
			DriftStdDevMS:   30.0,
			PosterURL:       makePosterURL("Testronic Labs", "dc2626", "ffffff"),
		},
		{
			ID:              "vnd-pixelogic",
			Name:            "Pixelogic Media",
			Components:      []string{"SUBTITLE"},
			Markets:         allFiveMarkets,
			HourlyRateUSD:   80.0,
			TurnaroundHours: 8.0,
			TargetAccuracy:  0.96,
			DriftMeanMS:     0.0,
			DriftStdDevMS:   0.0,
			PosterURL:       makePosterURL("Pixelogic Media", "059669", "ffffff"),
		},
		{
			ID:              "vnd-sound-vision-india",
			Name:            "Sound & Vision India",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         indianMarkets,
			HourlyRateUSD:   90.0,
			TurnaroundHours: 20.0,
			TargetAccuracy:  0.95,
			DriftMeanMS:     35.0,
			DriftStdDevMS:   15.0,
			PosterURL:       makePosterURL("Sound Vision", "f59e0b", "ffffff"),
		},
		{
			ID:              "vnd-prasad",
			Name:            "Prasad Corp",
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         indianMarkets,
			HourlyRateUSD:   110.0,
			TurnaroundHours: 16.0,
			TargetAccuracy:  0.92,
			DriftMeanMS:     40.0,
			DriftStdDevMS:   18.0,
			PosterURL:       makePosterURL("Prasad Corp", "10b981", "ffffff"),
		},
		// Technicolor: Dedicated Video/QC facility.
		// Coverage semantics: VIDEO is global and market-agnostic (Markets is empty).
		{
			ID:              "vnd-technicolor",
			Name:            "Technicolor",
			Components:      []string{"VIDEO"},
			Markets:         []string{},
			HourlyRateUSD:   185.0,
			TurnaroundHours: 16.0,
			TargetAccuracy:  0.98,
			DriftMeanMS:     0.0,
			DriftStdDevMS:   0.0,
			PosterURL:       makePosterURL("Technicolor", "d97706", "ffffff"),
		},
		// Prime Focus: Video/QC & Subtitling lab.
		// Coverage semantics: VIDEO is global and market-agnostic. Markets (all 5) gate its SUBTITLE coverage only.
		{
			ID:              "vnd-prime-focus",
			Name:            "Prime Focus",
			Components:      []string{"VIDEO", "SUBTITLE"},
			Markets:         allFiveMarkets,
			HourlyRateUSD:   75.0,
			TurnaroundHours: 28.0,
			TargetAccuracy:  0.89,
			DriftMeanMS:     0.0,
			DriftStdDevMS:   0.0,
			PosterURL:       makePosterURL("Prime Focus", "ea580c", "ffffff"),
		},
	}
}

// GenerateFillerVendors returns additional named facilities up to count for large scale testing.
func GenerateFillerVendors(count int, rng *types.RNG) []VendorSpec {
	if count <= 0 {
		return nil
	}

	result := make([]VendorSpec, 0, count)
	westernMarkets := []string{"en-US", "de-DE", "fr-FR"}

	for i := 0; i < count; i++ {
		cityName := "Munich"
		if rng != nil && rng.Faker != nil {
			cityName = rng.Faker.City()
		}

		name := fmt.Sprintf("%s Localization Labs", cityName)
		slug := fmt.Sprintf("vnd-%s-%03d", strings.ToLower(strings.ReplaceAll(cityName, " ", "-")), i+1)
		rate := 75.0 + float64((i*15)%100)
		turnaround := 16.0 + float64((i*8)%32)
		acc := 0.88 + float64((i*3)%10)*0.01

		result = append(result, VendorSpec{
			ID:              slug,
			Name:            name,
			Components:      []string{"AUDIO", "SUBTITLE"},
			Markets:         westernMarkets,
			HourlyRateUSD:   rate,
			TurnaroundHours: turnaround,
			TargetAccuracy:  acc,
			DriftMeanMS:     45.0,
			DriftStdDevMS:   20.0,
			PosterURL:       makePosterURL(name, "334155", "ffffff"),
		})
	}

	return result
}
