package entities

import "github.com/elliot14A/fincher/pkg/domain/models"

// World encapsulates the starter graph of SQLite domain entities (Vendors + Titles) assembled for database seeding.
type World struct {
	Vendors []*models.Vendor
	Titles  []*models.Title
}
