package aoclib

import "testing"

func TestSessionCookie(t *testing.T) {
	helper := NewAoCHelper()
	if helper.sessionCookie == "" {
		t.Fatalf("AoCHelper doesnt have a session cookie.")
	}
}

func TestInputWithoutCache(t *testing.T) {
	helper := NewAoCHelper()
	input := helper.GetInput(2023, 01, true)
	if input[0] != "76xkqjzqtwonfour" {
		t.Fatalf("Retrieved input does not match expected input")
	}
}
