package homehub

import "testing"

func TestNotAuthenticated(t *testing.T) {
	a := &authData{}

	if a.isAuthenticated() != false {
		t.Errorf("Expected isAuthenticated to be false")
	}
}

func TestAuthenticated(t *testing.T) {
	a := &authData{"http://foo/bar", "john", "doe", "87987", "222121", 0}

	if a.isAuthenticated() != true {
		t.Errorf("Expected isAuthenticated to be true")
	}
}
