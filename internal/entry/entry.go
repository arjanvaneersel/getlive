package entry

import (
	"context"
	"database/sql"
	"time"

	"github.com/arjanvaneersel/getlive/internal/platform/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

const collection = "entries"

var (
	// ErrNotFound is used when a specific Entry is requested but does not exist.
	ErrNotFound = errors.New("Entry not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// List retrieves a list of existing entries from the database.
func List(ctx context.Context, db *sqlx.DB) ([]Entry, error) {
	ctx, span := trace.StartSpan(ctx, "internal.entry.List")
	defer span.End()

	r := []Entry{}
	const q = `SELECT * FROM entries`

	if err := db.SelectContext(ctx, &r, q); err != nil {
		return nil, errors.Wrap(err, "selecting entries")
	}

	return r, nil
}

// Retrieve gets the specified entry from the database.
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Entry, error) {
	ctx, span := trace.StartSpan(ctx, "internal.entry.Retrieve")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	var r Entry
	const q = `SELECT * FROM entries WHERE entry_id = $1`
	if err := db.GetContext(ctx, &r, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, errors.Wrapf(err, "selecting entry %q", id)
	}

	return &r, nil
}

// Create inserts a new entry into the database.
func Create(ctx context.Context, db *sqlx.DB, n NewEntry, now time.Time) (*Entry, error) {
	ctx, span := trace.StartSpan(ctx, "internal.entry.Create")
	defer span.End()

	r := Entry{
		ID:               uuid.New().String(),
		Time:             n.Time,
		Title:            n.Title,
		Description:      n.Description,
		URL:              n.URL,
		Categories:       n.Categories,
		Keywords:         n.Keywords,
		SocialmediaLinks: n.SocialmediaLinks,
		Approved:         n.Approved,
		ApprovedBy:       n.ApprovedBy,
		Owner:            n.Owner,
		DateCreated:      now.UTC(),
		DateUpdated:      now.UTC(),
	}

	const q = `INSERT INTO entries
		(entry_id, date_time, title, description, url, categories, keywords, 
		socialmedia_links, approved, approved_by, owner, date_created, 
		date_updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := db.ExecContext(
		ctx, q,
		r.ID, r.Time, r.Title, r.Description,
		r.URL, r.Categories, r.Keywords, r.SocialmediaLinks,
		r.Approved, r.ApprovedBy, r.Owner,
		r.DateCreated, r.DateUpdated,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting entry")
	}

	return &r, nil
}

// Update replaces an entry document in the database.
func Update(ctx context.Context, db *sqlx.DB, claims auth.Claims, id string, upd UpdateEntry, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "internal.entry.Update")
	defer span.End()

	r, err := Retrieve(ctx, db, id)
	if err != nil {
		return err
	}

	if upd.Time != nil {
		r.Time = *upd.Time
	}
	if upd.Title != nil {
		r.Title = *upd.Title
	}
	if upd.Description != nil {
		r.Description = *upd.Description
	}
	if upd.URL != nil {
		r.URL = *upd.URL
	}
	if upd.Categories != nil {
		r.Categories = upd.Categories
	}
	if upd.SocialmediaLinks != nil {
		r.SocialmediaLinks = upd.SocialmediaLinks
	}
	if upd.Keywords != nil {
		r.Keywords = upd.Keywords
	}
	if upd.Approved != nil {
		r.Approved = *upd.Approved
	}
	if upd.ApprovedBy != nil {
		r.ApprovedBy = *upd.ApprovedBy
	}
	if upd.Owner != nil {
		r.Owner = upd.Owner
	}

	r.DateUpdated = now

	const q = `UPDATE entries SET
		"date_time" = $2,
		"title" = $3,
		"description" = $4,
		"url" = $5,
		"categories" = $6,
		"keywords" = $7,
		"socialmedia_links" = $8,
		"approved" = $9,
		"approved_by" = $10,
		"owner" = $11,
		"date_updated" = $12
		WHERE entry_id = $1`
	_, err = db.ExecContext(ctx, q, id,
		r.Time, r.Title, r.Description,
		r.URL, r.Categories, r.Keywords,
		r.SocialmediaLinks, r.Approved,
		r.ApprovedBy, r.Owner, r.DateUpdated,
	)
	if err != nil {
		return errors.Wrap(err, "updating entry")
	}

	return nil
}

// Delete removes a entry from the database.
func Delete(ctx context.Context, db *sqlx.DB, id string) error {
	ctx, span := trace.StartSpan(ctx, "internal.entry.Delete")
	defer span.End()

	if _, err := uuid.Parse(id); err != nil {
		return ErrInvalidID
	}

	const q = `DELETE FROM entries WHERE entry_id = $1`

	if _, err := db.ExecContext(ctx, q, id); err != nil {
		return errors.Wrapf(err, "deleting entry %s", id)
	}

	return nil
}
