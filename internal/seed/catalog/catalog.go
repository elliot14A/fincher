package catalog

import "github.com/elliot14A/fincher/pkg/domain/models"

// MarketLocale defines a target territory market and its localized language.
type MarketLocale struct {
	Tag      string // "de-DE" — the market key
	Language string // "de-DE" — same as Tag
	Country  string // "DE" ("IN" for both hi-IN and te-IN)
	Name     string // "German"
}

// Markets returns the 5 canonical language-market locales.
func Markets() []MarketLocale {
	return []MarketLocale{
		{Tag: "en-US", Language: "en-US", Country: "US", Name: "English"},
		{Tag: "de-DE", Language: "de-DE", Country: "DE", Name: "German"},
		{Tag: "fr-FR", Language: "fr-FR", Country: "FR", Name: "French"},
		{Tag: "hi-IN", Language: "hi-IN", Country: "IN", Name: "Hindi"},
		{Tag: "te-IN", Language: "te-IN", Country: "IN", Name: "Telugu"},
	}
}

// ComponentTypes returns the standard localization components.
func ComponentTypes() []models.ComponentType {
	return []models.ComponentType{
		models.ComponentAudio,
		models.ComponentSubtitle,
		models.ComponentVideo,
	}
}
