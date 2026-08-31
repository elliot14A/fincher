package catalog

import "github.com/elliot14A/fincher/pkg/domain/models"

// TitleSpec defines the configuration profile for a curated title.
type TitleSpec struct {
	ID                  string
	Name                string
	Slug                string
	Type                models.TitleType
	PremiereOffsetHours float64
	PosterURL           string
	Genre               string
	Synopsis            string
	OverallStatus       models.TitleStatus
	IsHero              bool
	Markets             []string // targeted language-market tags
}

// CuratedTitles returns the list of curated real-world titles with verified poster image URLs and market footprints.
func CuratedTitles() []TitleSpec {
	allFiveMarkets := []string{"en-US", "de-DE", "fr-FR", "hi-IN", "te-IN"}
	westernMarkets := []string{"en-US", "de-DE", "fr-FR"}
	indianFirstMarkets := []string{"te-IN", "hi-IN", "en-US"}

	return []TitleSpec{
		{
			ID:                  "title-avatar-fire-ash",
			Name:                "Avatar: Fire and Ash",
			Slug:                "avatar-fire-ash",
			Type:                models.TitleTypeFeature,
			PremiereOffsetHours: 72.0,
			PosterURL:           "https://upload.wikimedia.org/wikipedia/en/9/95/Avatar_Fire_and_Ash_poster.jpeg",
			Genre:               "Sci-Fi / Adventure",
			Synopsis:            "Jake Sully and Neytiri encounter a new, aggressive volcanic clan of Na'vi on Pandora.",
			OverallStatus:       models.StatusOnTrack,
			IsHero:              true,
			Markets:             allFiveMarkets,
		},
		{
			ID:                  "title-pushpa-the-rule",
			Name:                "Pushpa 2: The Rule",
			Slug:                "pushpa-the-rule",
			Type:                models.TitleTypeFeature,
			PremiereOffsetHours: 120.0,
			PosterURL:           "https://image.tmdb.org/t/p/w780/4FwM8TqXG2vYy9x1F9l5aM0mO1k.jpg",
			Genre:               "Action / Thriller",
			Synopsis:            "Pushpa Raj consolidates his red sanders smuggling syndicate while evading high-ranking police forces.",
			OverallStatus:       models.StatusProcessing,
			IsHero:              false,
			Markets:             indianFirstMarkets,
		},
		{
			ID:                  "title-dune-part-two",
			Name:                "Dune: Part Two",
			Slug:                "dune-part-two",
			Type:                models.TitleTypeFeature,
			PremiereOffsetHours: 360.0,
			PosterURL:           "https://image.tmdb.org/t/p/w780/1pdfLvkbY9ohJlCjQH2CZjjYVvJ.jpg",
			Genre:               "Sci-Fi / Drama",
			Synopsis:            "Paul Atreides unites with Chani and the Fremen while seeking revenge against the conspirators who destroyed his family.",
			OverallStatus:       models.StatusProcessing,
			IsHero:              false,
			Markets:             westernMarkets,
		},
		{
			ID:                  "title-shogun",
			Name:                "Shōgun",
			Slug:                "shogun",
			Type:                models.TitleTypeSeries,
			PremiereOffsetHours: 240.0,
			PosterURL:           "https://image.tmdb.org/t/p/w780/7O4iVfOMQmdCSxhOg1WnzG1AgYT.jpg",
			Genre:               "Drama / Historical",
			Synopsis:            "Lord Yoshii Toranaga discovers secrets that could tip the scales of power in feudal Japan.",
			OverallStatus:       models.StatusOnTrack,
			IsHero:              false,
			Markets:             westernMarkets,
		},
		{
			ID:                  "title-the-batman-part-ii",
			Name:                "The Batman Part II",
			Slug:                "the-batman-part-ii",
			Type:                models.TitleTypeFeature,
			PremiereOffsetHours: 480.0,
			PosterURL:           "https://image.tmdb.org/t/p/w780/74xTEgt7R36Fpooo50r9T25onhq.jpg",
			Genre:               "Action / Crime",
			Synopsis:            "Bruce Wayne delves deeper into the corruption of Gotham City's criminal underworld.",
			OverallStatus:       models.StatusOnTrack,
			IsHero:              false,
			Markets:             allFiveMarkets,
		},
		{
			ID:                  "title-deadpool-and-wolverine",
			Name:                "Deadpool & Wolverine",
			Slug:                "deadpool-and-wolverine",
			Type:                models.TitleTypeFeature,
			PremiereOffsetHours: -720.0, // Historical premiere (already shipped)
			PosterURL:           "https://image.tmdb.org/t/p/w780/8cdWjvZQUExUUTzyp4t6EDMubfO.jpg",
			Genre:               "Action / Comedy",
			Synopsis:            "Wade Wilson teams up with a reluctant Wolverine to save the multiverse.",
			OverallStatus:       models.StatusShipped,
			IsHero:              false,
			Markets:             allFiveMarkets,
		},
		{
			ID:                  "title-oppenheimer",
			Name:                "Oppenheimer",
			Slug:                "oppenheimer",
			Type:                models.TitleTypeFeature,
			PremiereOffsetHours: -1440.0, // Historical premiere (already shipped)
			PosterURL:           "https://image.tmdb.org/t/p/w780/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
			Genre:               "Drama / History",
			Synopsis:            "The story of American scientist J. Robert Oppenheimer and his role in the development of the atomic bomb.",
			OverallStatus:       models.StatusShipped,
			IsHero:              false,
			Markets:             westernMarkets,
		},
	}
}
