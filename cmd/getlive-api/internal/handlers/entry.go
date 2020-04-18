package handlers

import (
	"context"
	"net/http"

	"github.com/arjanvaneersel/getlive/internal/entry"
	"github.com/arjanvaneersel/getlive/internal/platform/auth"
	"github.com/arjanvaneersel/getlive/internal/platform/web"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Entry represents the Entry API method handler set.
type Entry struct {
	db *sqlx.DB

	// ADD OTHER STATE LIKE THE LOGGER IF NEEDED.
}

// List gets all existing entries in the system.
func (p *Entry) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Entry.List")
	defer span.End()

	o, err := entry.List(ctx, p.db)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, o, http.StatusOK)
}

// Retrieve returns the specified entry from the system.
func (p *Entry) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Entry.Retrieve")
	defer span.End()

	o, err := entry.Retrieve(ctx, p.db, params["id"])
	if err != nil {
		switch err {
		case entry.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case entry.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, o, http.StatusOK)
}

// Create decodes the body of a request to create a new entry. The full
// entry with generated fields is sent back in the response.
func (p *Entry) Create(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Entry.Create")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var n entry.NewEntry
	if err := web.Decode(r, &n); err != nil {
		return errors.Wrap(err, "decoding new entry")
	}

	o, err := entry.Create(ctx, p.db, claims, n, v.Now)
	if err != nil {
		return errors.Wrapf(err, "creating new entry: %+v", n)
	}

	return web.Respond(ctx, w, o, http.StatusCreated)
}

// Update decodes the body of a request to update an existing entry. The ID
// of the entry is part of the request URL.
func (p *Entry) Update(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Entry.Update")
	defer span.End()

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return web.NewShutdownError("claims missing from context")
	}

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var up entry.UpdateEntry
	if err := web.Decode(r, &up); err != nil {
		return errors.Wrap(err, "")
	}

	if err := entry.Update(ctx, p.db, claims, params["id"], up, v.Now); err != nil {
		switch err {
		case entry.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case entry.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case entry.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating entry %q: %+v", params["id"], up)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// Delete removes a single entry identified by an ID in the request URL.
func (p *Entry) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Entry.Delete")
	defer span.End()

	if err := entry.Delete(ctx, p.db, params["id"]); err != nil {
		switch err {
		case entry.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "Id: %s", params["id"])
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
