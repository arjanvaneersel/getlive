package entry

import (
	"time"

	"github.com/lib/pq"
)

// Entry represents an entry into the system.
type Entry struct {
	ID               string         `db:"entry_id" json:"id"`
	Time             time.Time      `db:"date_time" json:"date_time"`
	Title            string         `db:"title" json:"title"`
	Description      string         `db:"description" json:"description"`
	URL              string         `db:"url" json:"url"`
	Categories       pq.StringArray `db:"categories" json:"categories"`
	Keywords         pq.StringArray `db:"keywords" json:"keywords"`
	SocialmediaLinks pq.StringArray `db:"socialmedia_links" json:"socialmedia_links"`
	Approved         bool           `db:"approved" json:"approved"`
	ApprovedBy       string         `db:"approved_by" json:"approved_by"`
	Owner            *string        `db:"owner" json:"owner"`
	DateCreated      time.Time      `db:"date_created" json:"date_created"`
	DateUpdated      time.Time      `db:"date_updated" json:"date_updated"`
}

// NewEntry contains information needed to create a new Entry.
type NewEntry struct {
	Time             time.Time `json:"date_time" validate:"required"`
	Title            string    `json:"title" validate:"required"`
	Description      string    `json:"description"`
	URL              string    `json:"url" validate:"required"`
	Categories       []string  `json:"categories"`
	Keywords         []string  `json:"keywords"`
	SocialmediaLinks []string  `json:"socialmedia_links"`

	// These fields can only be set by admins.
	// If a non-admin posts a new entry, approval will be
	// set to false and the owner will be the user posting.
	Approved   bool    `json:"approved"`
	ApprovedBy string  `json:"approved_by"`
	Owner      *string `json:"owner"`
}

// UpdateEntry defines what information may be provided to modify an existing
// Entry. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateEntry struct {
	Time             *time.Time `json:"date_time"`
	Title            *string    `json:"title"`
	Description      *string    `json:"description"`
	URL              *string    `json:"url"`
	Categories       []string   `json:"categories"`
	Keywords         []string   `json:"keywords"`
	SocialmediaLinks []string   `json:"socialmedia_links"`

	// These fields can only be set by admins.
	// If a non-admin posts a new entry, approval will be
	// set to false and the owner will be the user posting.
	Approved   *bool   `json:"approved"`
	ApprovedBy *string `json:"approved_by"`
	Owner      *string `json:"owner"`
}
