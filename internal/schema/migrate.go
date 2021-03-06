package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed from this slice once they have been ran in
// production.
//
// Using constants in a .go file is an easy way to ensure the queries are part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add users",
		Script: `
CREATE TABLE users (
	user_id       UUID DEFAULT '00000000-0000-0000-0000-000000000000',
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,

	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     2,
		Description: "Add entries",
		Script: `
CREATE TABLE entries (
	entry_id          UUID,
	date_time         TIMESTAMP,
	title             TEXT,
	description       TEXT,
	url               TEXT,
	categories        TEXT[],
	keywords          TEXT[],
	socialmedia_links TEXT[],
	approved          BOOL,
	approved_by       TEXT,
	owner             UUID DEFAULT '00000000-0000-0000-0000-000000000000',
	date_created      TIMESTAMP,
	date_updated      TIMESTAMP,

	PRIMARY KEY (entry_id),
	FOREIGN KEY (owner) REFERENCES users(user_id) ON DELETE CASCADE
);`,
	},
}
