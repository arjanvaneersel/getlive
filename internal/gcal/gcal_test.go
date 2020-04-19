package gcal

import (
	"testing"
)

// func TestNew(t *testing.T) {
// 	_, err := New("../../credentials.json")
// 	if err != nil {
// 		t.Fatalf("expected to pass, but got: %v", err)
// 	}
// }

func TestPost(t *testing.T) {
	c, err := New("../../credentials.json")
	if err != nil {
		t.Fatalf("expected to pass, but got: %v", err)
	}

	c.Post()
}

func TestList(t *testing.T) {
	c, err := New("../../credentials.json")
	if err != nil {
		t.Fatalf("expected to pass, but got: %v", err)
	}

	c.List()
}
