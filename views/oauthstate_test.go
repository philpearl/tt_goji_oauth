package views

import (
	"testing"
)

func TestOauthState(t *testing.T) {
	s := newOauthState()
	s.ProviderName = "dummy"

	state := s.Encode()

	s2, err := newOauthStateFromString(state)
	if err != nil {
		t.Fatalf("could not decode oauthstate from string. %v", err)
	}

	if s.ProviderName != s2.ProviderName {
		t.Fatalf("provider names differ %s %s", s.ProviderName, s2.ProviderName)
	}

	if s.Secret != s2.Secret {
		t.Fatalf("secrets differ %s %s", s.Secret, s2.Secret)
	}
}
