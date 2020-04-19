package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/arjanvaneersel/getlive/internal/entry"
	"github.com/arjanvaneersel/getlive/internal/gcal"
	"github.com/arjanvaneersel/getlive/internal/platform/auth"
	"github.com/arjanvaneersel/getlive/internal/platform/web"
	"github.com/arjanvaneersel/getlive/internal/user"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type WebAdmin struct {
	db     *sqlx.DB
	claims *auth.Claims
	cal    *gcal.Calendar
}

// List returns all the existing users in the system.
func (wa *WebAdmin) List(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.List")
	defer span.End()

	entries, err := entry.List(ctx, wa.db)
	if err != nil {
		return err
	}

	var t string
	if wa.claims == nil {
		t = loginTemplate
	} else {
		t = listTemplate
	}

	tpl, err := template.New("").Parse(t)
	if err != nil {
		return err
	}

	return tpl.Execute(w, entries)
}

func (wa *WebAdmin) Login(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.Login")
	defer span.End()

	r.ParseForm()

	claims, err := user.Authenticate(ctx, wa.db, time.Now(), r.Form.Get("email"), r.Form.Get("password"))
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating")
		}
	}

	wa.claims = &claims
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func (wa *WebAdmin) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.Retrieve")
	defer span.End()

	if wa.claims == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}

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

func (wa *WebAdmin) Approve(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.Retrieve")
	defer span.End()

	if wa.claims == nil {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}

	yes := true
	r.ParseForm()
	up := entry.UpdateEntry{
		Approved:   &yes,
		ApprovedBy: &wa.claims.Subject,
	}

	if len(r.Form.Get("categories")) > 0 {
		//TODO: Check that categories are in Concerts & Festivals, MIND, BODY, SOUL, Covid19
		categories := strings.Split(r.Form.Get("categories"), " ")
		up.Categories = categories
	}

	if err := entry.Update(ctx, wa.db, *wa.claims, params["id"], up, time.Now()); err != nil {
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

	// Get all data for publishing to the calendar
	// if the call service is initialized
	if wa.cal != nil {
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

		title := e.Title
		if len(e.Categories) > 0 {
			title = fmt.Sprintf("[%s] %s", strings.ToUpper(e.Categories[0]), title)
		}

		// TODO: Store cal link in DB
		// TODO: Better time handling
		if _, err := wa.cal.Post(title, e.Description, e.URL, e.Time, e.Time.Add(1*time.Hour)); err != nil {
			return errors.Wrap(err, "publishing to calendar")
		}
	}

	http.Redirect(w, r, "/entries/"+params["id"], http.StatusFound)
	return nil
}

func (wa *WebAdmin) OAuthCallback(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	ctx, span := trace.StartSpan(ctx, "handlers.WebAdmin.OAuthCallback")
	defer span.End()

	fmt.Fprintf(w, "Code: %s\n", r.FormValue("code"))

	return web.Respond(ctx, w, nil, http.StatusOK)
}
