package entry_test

import (
	"context"
	"testing"
	"time"

	"github.com/arjanvaneersel/getlive/internal/entry"
	"github.com/arjanvaneersel/getlive/internal/platform/auth"
	"github.com/arjanvaneersel/getlive/internal/tests"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
)

// TestEntry validates the full set of CRUD operations on Entry values.
func TestEntry(t *testing.T) {
	db, teardown := tests.NewUnit(t)
	defer teardown()

	t.Log("Given the need to work with Entry records.")
	{
		t.Log("\tWhen handling a single Entry.")
		{
			now := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)
			n := entry.NewEntry{
				Time:             now,
				Title:            "Test event",
				Description:      "This is a test event.",
				URL:              "http://example.com",
				SocialmediaLinks: []string{"http://instagram.com/test"},
				Approved:         false,
			}

			ctx := context.Background()

			claims := auth.NewClaims(
				"718ffbea-f4a1-4667-8ae3-b349da52675e", // This is just some random UUID.
				[]string{auth.RoleAdmin, auth.RoleUser},
				now, time.Hour,
			)

			r, err := entry.Create(ctx, db, claims, n, now)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create an entry : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to create an entry.", tests.Success)

			saved, err := entry.Retrieve(ctx, db, r.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve entry by ID: %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve entry by ID.", tests.Success)

			if diff := cmp.Diff(r, saved); diff != "" {
				t.Fatalf("\t%s\tShould get back the same entry. Diff:\n%s", tests.Failed, diff)
			}
			t.Logf("\t%s\tShould get back the same entry.", tests.Success)

			upd := entry.UpdateEntry{
				Title: tests.StringPointer("Testing event"),
			}
			updatedTime := time.Date(2019, time.January, 1, 1, 1, 1, 0, time.UTC)

			if err := entry.Update(ctx, db, claims, r.ID, upd, updatedTime); err != nil {
				t.Fatalf("\t%s\tShould be able to update entry: %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update entry.", tests.Success)

			saved, err = entry.Retrieve(ctx, db, r.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve updated entry : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve updated entry.", tests.Success)

			// Check specified fields were updated. Make a copy of the original entry
			// and change just the fields we expect then diff it with what was saved.
			want := *r
			want.Title = *upd.Title
			want.DateUpdated = updatedTime

			if diff := cmp.Diff(want, *saved); diff != "" {
				t.Fatalf("\t%s\tShould get back the same entry. Diff:\n%s", tests.Failed, diff)
			}
			t.Logf("\t%s\tShould get back the same entry.", tests.Success)

			upd = entry.UpdateEntry{
				Description: tests.StringPointer("This is a testing event"),
			}

			if err := entry.Update(ctx, db, claims, r.ID, upd, updatedTime); err != nil {
				t.Fatalf("\t%s\tShould be able to update just some fields of entry : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to update just some fields of entry.", tests.Success)

			saved, err = entry.Retrieve(ctx, db, r.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve updated entry : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve updated entry.", tests.Success)

			if saved.Description != *upd.Description {
				t.Fatalf("\t%s\tShould be able to see updated Description field : got %q want %q.", tests.Failed, saved.Description, *upd.Description)
			} else {
				t.Logf("\t%s\tShould be able to see updated Description field.", tests.Success)
			}

			if err := entry.Delete(ctx, db, r.ID); err != nil {
				t.Fatalf("\t%s\tShould be able to delete entry : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to delete entry.", tests.Success)

			saved, err = entry.Retrieve(ctx, db, r.ID)
			if errors.Cause(err) != entry.ErrNotFound {
				t.Fatalf("\t%s\tShould NOT be able to retrieve deleted entry : %s.", tests.Failed, err)
			}
			t.Logf("\t%s\tShould NOT be able to retrieve deleted entry.", tests.Success)
		}
	}
}
