package aoclib

import "testing"

func TestSessionCookie(t *testing.T) {
	helper := NewAoCHelper(2023, 01)
	if helper.sessionCookie == "" {
		t.Fatalf("AoCHelper doesnt have a session cookie.")
	}
}

func TestInputWithoutCache(t *testing.T) {
	helper := NewAoCHelper(2023, 01)
	input := helper.GetInput(true)
	if input[0] != "76xkqjzqtwonfour" {
		t.Fatalf("Retrieved input does not match expected input")
	}
}

func TestSubmit(t *testing.T) {
	helper := NewAoCHelper(2023, 01)
	helper.Submit(1, 54644)
}
