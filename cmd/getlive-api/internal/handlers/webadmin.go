package handlers

import (
	"context"
	"github.com/arjanvaneersel/getlive/internal/entry"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"html/template"
	"net/http"
)

type WebAdmin struct {
	db  *sqlx.DB
	jwt string
}

// List returns all the existing users in the system.
func (wa *WebAdmin) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.List")
	defer span.End()

	entries, err := entry.List(ctx, wa.db)
	if err != nil {
		return err
	}

	tpl, err := template.New("").Parse(listTemplate)
	if err != nil {
		return err
	}

	return tpl.Execute(w, entries)
}

func (wa *WebAdmin) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.Retrieve")
	defer span.End()

	// claims, ok := ctx.Value(auth.Key).(auth.Claims)
	// if !ok {
	// 	return errors.New("claims missing from context")
	// }

	e, err := entry.Retrieve(ctx, wa.db, params["id"])
	if err != nil {
		switch err {
		case entry.ErrInvalidID:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return nil
		case entry.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return nil
		case entry.ErrForbidden:
			http.Error(w, err.Error(), http.StatusForbidden)
			return nil
		default:
			return errors.Wrapf(err, "Id: %s", params["id"])
		}
	}

	tpl, err := template.New("").Parse(entryTemplate)
	if err != nil {
		return err
	}

	return tpl.Execute(w, e)
}
