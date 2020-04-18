package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/arjanvaneersel/getlive/internal/mid"
	"github.com/arjanvaneersel/getlive/internal/platform/auth" // Import is removed in final PR
	"github.com/arjanvaneersel/getlive/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// API constructs an http.Handler with all application routes defined.
func API(build string, shutdown chan os.Signal, log *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log), mid.Metrics(), mid.Panics(log))

	// Register health check endpoint. This route is not authenticated.
	check := Check{
		build: build,
		db:    db,
	}
	app.Handle("GET", "/v1/health", check.Health)

	// Register user management and authentication endpoints.
	u := User{
		db:            db,
		authenticator: authenticator,
	}
	app.Handle("GET", "/v1/users", u.List, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("POST", "/v1/users", u.Create, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("GET", "/v1/users/:id", u.Retrieve, mid.Authenticate(authenticator))
	app.Handle("PUT", "/v1/users/:id", u.Update, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
	app.Handle("DELETE", "/v1/users/:id", u.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// This route is not authenticated
	app.Handle("GET", "/v1/users/token", u.Token)

	// Register entry management endpoints.
	e := Entry{
		db: db,
	}
	app.Handle("GET", "/v1/entries", e.List)
	app.Handle("POST", "/v1/entries", e.Create, mid.Authenticate(authenticator))
	app.Handle("GET", "/v1/entries/:id", e.Retrieve)
	app.Handle("PUT", "/v1/entries/:id", e.Update, mid.Authenticate(authenticator))
	app.Handle("DELETE", "/v1/entries/:id", e.Delete, mid.Authenticate(authenticator))

	w := WebAdmin{
		db: db,
	}

	app.Handle("GET", "/", w.List)
	app.Handle("GET", "/entries/:id", w.Retrieve)

	return app
}
