package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arjanvaneersel/getlive/cmd/getlive-api/internal/handlers"
	"github.com/arjanvaneersel/getlive/internal/entry"
	"github.com/arjanvaneersel/getlive/internal/platform/web"
	"github.com/arjanvaneersel/getlive/internal/tests"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TestEntries runs a series of tests to exercise Entry behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestEntries(t *testing.T) {
	test := tests.NewIntegration(t)
	defer test.Teardown()

	shutdown := make(chan os.Signal, 1)
	tests := EntryTests{
		app:       handlers.API("develop", shutdown, test.Log, test.DB, test.Authenticator),
		userToken: test.Token("admin@example.com", "gophers"),
	}

	t.Run("postEntry400", tests.postEntry400)
	t.Run("postEntry401", tests.postEntry401)
	t.Run("getEntry404", tests.getEntry404)
	t.Run("getEntry400", tests.getEntry400)
	t.Run("deleteEntryNotFound", tests.deleteEntryNotFound)
	t.Run("putEntry404", tests.putEntry404)
	t.Run("crudEntries", tests.crudEntry)
}

// EntryTests holds methods for each entry subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type EntryTests struct {
	app       http.Handler
	userToken string
}

// postEntry400 validates an entry can't be created with the endpoint
// unless a valid entry document is submitted.
func (pt *EntryTests) postEntry400(t *testing.T) {
	r := httptest.NewRequest("POST", "/v1/entries", strings.NewReader(`{}`))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new entry can't be created with an invalid document.")
	{
		t.Log("\tTest 0:\tWhen using an incomplete entry value.")
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tShould receive a status code of 400 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 400 for the response.", tests.Success)

			// Inspect the response.
			var got web.ErrorResponse
			if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response to an error type : %v", tests.Failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the response to an error type.", tests.Success)

			// Define what we want to see.
			want := web.ErrorResponse{
				Error: "field validation error",
				Fields: []web.FieldError{
					{Field: "date_time", Error: "date_time is a required field"},
					{Field: "title", Error: "title is a required field"},
					{Field: "description", Error: "description is a required field"},
					{Field: "url", Error: "url is a required field"},
				},
			}

			// We can't rely on the order of the field errors so they have to be
			// sorted. Tell the cmp package how to sort them.
			sorter := cmpopts.SortSlices(func(a, b web.FieldError) bool {
				return a.Field < b.Field
			})

			if diff := cmp.Diff(want, got, sorter); diff != "" {
				t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// postEntry401 validates an entry can't be created with the endpoint
// unless the user is authenticated
func (pt *EntryTests) postEntry401(t *testing.T) {
	n := entry.NewEntry{
		Time:        time.Now(),
		Title:       "Test event",
		Description: "This is a testing event",
		URL:         "http://example.com",
	}

	body, err := json.Marshal(&n)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/v1/entries", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// Not setting an authorization header

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate a new entry can't be created with an invalid document.")
	{
		t.Log("\tTest 0:\tWhen using an incomplete entry value.")
		{
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("\t%s\tShould receive a status code of 401 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 401 for the response.", tests.Success)
		}
	}
}

// getEntry400 validates an entry request for a malformed id.
func (pt *EntryTests) getEntry400(t *testing.T) {
	id := "12345"

	r := httptest.NewRequest("GET", "/v1/entries/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting an entry with a malformed id.")
	{
		t.Logf("\tTest 0:\tWhen using the new entry %s.", id)
		{
			if w.Code != http.StatusBadRequest {
				t.Fatalf("\t%s\tShould receive a status code of 400 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 400 for the response.", tests.Success)

			recv := w.Body.String()
			resp := `{"error":"ID is not in its proper form"}`
			if resp != recv {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// getEntry404 validates an entry request for an entry that does not exist with the endpoint.
func (pt *EntryTests) getEntry404(t *testing.T) {
	id := "a224a8d6-3f9e-4b11-9900-e81a25d80702"

	r := httptest.NewRequest("GET", "/v1/entries/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting an entry with an unknown id.")
	{
		t.Logf("\tTest 0:\tWhen using the new entry %s.", id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tShould receive a status code of 404 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 404 for the response.", tests.Success)

			recv := w.Body.String()
			resp := "Entry not found"
			if !strings.Contains(recv, resp) {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// deleteEntryNotFound validates deleting an entry that does not exist is not a failure.
func (pt *EntryTests) deleteEntryNotFound(t *testing.T) {
	id := "112262f1-1a77-4374-9f22-39e575aa6348"

	r := httptest.NewRequest("DELETE", "/v1/entries/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting an entry that does not exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new entry %s.", id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", tests.Success)
		}
	}
}

// putEntry404 validates updating an entry that does not exist.
func (pt *EntryTests) putEntry404(t *testing.T) {
	up := entry.UpdateEntry{
		Title: tests.StringPointer("Nonexistent"),
	}

	id := "9b468f90-1cf1-4377-b3fa-68b450d632a0"

	body, err := json.Marshal(&up)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("PUT", "/v1/entries/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate updating an entry that does not exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new entry %s.", id)
		{
			if w.Code != http.StatusNotFound {
				t.Fatalf("\t%s\tShould receive a status code of 404 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 404 for the response.", tests.Success)

			recv := w.Body.String()
			resp := "Entry not found"
			if !strings.Contains(recv, resp) {
				t.Log("Got :", recv)
				t.Log("Want:", resp)
				t.Fatalf("\t%s\tShould get the expected result.", tests.Failed)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// crudEntry performs a complete test of CRUD against the api.
func (pt *EntryTests) crudEntry(t *testing.T) {
	p := pt.postEntry201(t)
	defer pt.deleteEntry204(t, p.ID)

	pt.getEntry200(t, p.ID)
	pt.putEntry204(t, p.ID)
}

// postEntry201 validates an entry can be created with the endpoint.
func (pt *EntryTests) postEntry201(t *testing.T) entry.Entry {
	eventDate, _ := time.Parse("2006-01-02", "2020-01-29")
	np := entry.NewEntry{
		Time:        eventDate,
		Title:       "Test event",
		Description: "This is a test event",
		URL:         "http://example.com",
	}

	body, err := json.Marshal(&np)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest("POST", "/v1/entries", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	// o is the value we will return.
	var o entry.Entry

	t.Log("Given the need to create a new entry with the entry endpoint.")
	{
		t.Log("\tTest 0:\tWhen using the declared entry value.")
		{
			if w.Code != http.StatusCreated {
				t.Fatalf("\t%s\tShould receive a status code of 201 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 201 for the response.", tests.Success)

			if err := json.NewDecoder(w.Body).Decode(&o); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like ID and Dates so we copy p.
			want := o
			want.Title = "Test event"
			want.Time = eventDate
			want.Description = "This is a test event"
			want.URL = "http://example.com"

			if diff := cmp.Diff(want, o); diff != "" {
				t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}

	return o
}

// deleteEntry200 validates deleting an entry that does exist.
func (pt *EntryTests) deleteEntry204(t *testing.T, id string) {
	r := httptest.NewRequest("DELETE", "/v1/entries/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate deleting an entry that does exist.")
	{
		t.Logf("\tTest 0:\tWhen using the new entry %s.", id)
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", tests.Success)
		}
	}
}

// getEntryEvent200 validates an entry request for an existing id.
func (pt *EntryTests) getEntry200(t *testing.T, id string) {
	r := httptest.NewRequest("GET", "/v1/entries/"+id, nil)
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to validate getting an entry that exists.")
	{
		t.Logf("\tTest 0:\tWhen using the new entry %s.", id)
		{
			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the response.", tests.Success)

			var o entry.Entry
			if err := json.NewDecoder(w.Body).Decode(&o); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
			}

			// Define what we wanted to receive. We will just trust the generated
			// fields like Dates so we copy p.
			want := o
			want.ID = id
			want.Title = "Test event"
			want.Description = "This is a test event"
			want.URL = "http://example.com"
			want.Time, _ = time.Parse("2006-01-02", "2020-01-29")

			if diff := cmp.Diff(want, o); diff != "" {
				t.Fatalf("\t%s\tShould get the expected result. Diff:\n%s", tests.Failed, diff)
			}
			t.Logf("\t%s\tShould get the expected result.", tests.Success)
		}
	}
}

// putEntry204 validates updating an entry that does exist.
func (pt *EntryTests) putEntry204(t *testing.T, id string) {
	body := `{"title": "Testing event"}`
	r := httptest.NewRequest("PUT", "/v1/entries/"+id, strings.NewReader(body))
	w := httptest.NewRecorder()

	r.Header.Set("Authorization", "Bearer "+pt.userToken)

	pt.app.ServeHTTP(w, r)

	t.Log("Given the need to update an entry with the entries endpoint.")
	{
		t.Log("\tTest 0:\tWhen using the modified entry value.")
		{
			if w.Code != http.StatusNoContent {
				t.Fatalf("\t%s\tShould receive a status code of 204 for the response : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 204 for the response.", tests.Success)

			r = httptest.NewRequest("GET", "/v1/entries/"+id, nil)
			w = httptest.NewRecorder()

			r.Header.Set("Authorization", "Bearer "+pt.userToken)

			pt.app.ServeHTTP(w, r)

			if w.Code != http.StatusOK {
				t.Fatalf("\t%s\tShould receive a status code of 200 for the retrieve : %v", tests.Failed, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of 200 for the retrieve.", tests.Success)

			var ro entry.Entry
			if err := json.NewDecoder(w.Body).Decode(&ro); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", tests.Failed, err)
			}

			if ro.Title != "Testing event" {
				t.Fatalf("\t%s\tShould see an updated Title : got %q want %q", tests.Failed, ro.Title, "Testing Event")
			}
			t.Logf("\t%s\tShould see an updated Title.", tests.Success)
		}
	}
}
