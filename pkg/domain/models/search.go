package models

// SearchResultKind enumerates the object kinds an operator can reference with @.
type SearchResultKind string

const (
	SearchKindTitle    SearchResultKind = "title"
	SearchKindVendor   SearchResultKind = "vendor"
	SearchKindDelivery SearchResultKind = "delivery"
)

// SearchResult is a unified, taggable reference to a domain object.
type SearchResult struct {
	Kind        SearchResultKind `json:"kind"`
	ID          string           `json:"id"`
	DisplayName string           `json:"display_name"`
	Identifier  string           `json:"identifier"`
	Subtitle    string           `json:"subtitle,omitempty"`
	ImageURL    string           `json:"image_url,omitempty"`
}
