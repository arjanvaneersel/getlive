package twitter

import (
	"os"
	"testing"
)

func TestGetNewEntryFromMedia(t *testing.T) {
	ytKey := os.Getenv("YT_KEY")
	if ytKey == "" {
		t.Log("ignoring test, because YT_KEY has not been provided in ENV")
		return
	}

	e, _, err := getNewEntryFromMedia("https://www.youtube.com/watch?v=bdOaMFD35C8", ytKey)
	if err != nil {
		t.Fatalf("expected to pass, but got: %v", err)
	}

	if e.Title == "" {
		t.Error("expected entry to have a title")
	}

	if e.Description == "" {
		t.Error("expected entry to have a description")
	}
	t.Logf("description: %s", e.Description)

	if e.URL == "" {
		t.Error("expected entry to have a URL")
	}
}
