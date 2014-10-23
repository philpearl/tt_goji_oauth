package providers

import (
	"github.com/golang/oauth2"
)

const (
	// username used by the provider
	PROVIDER_USERNAME = "username"
	// user's email address
	PROVIDER_EMAIL = "email"
	// User's (potentially) real name
	PROVIDER_NAME = "name"
	// id for the user - one that doesn't change
	PROVIDER_ID = "id"
)

type Provider interface {
	/*
	   GetName returns the name of the provider
	*/
	GetName() string
	/*
			GetUserInfo gets information about the user.

		    Information is returned in a map. Keys are defined by the PROVIDER_ constants
	*/
	GetUserInfo(t *oauth2.Transport) (map[string]interface{}, error)

	/*
	   GetConfig returns the oauth2 config for this provider
	*/
	GetConfig() *oauth2.Config
}

type GenericProvider struct {
	Name   string
	Config *oauth2.Config
}

func (p *GenericProvider) GetConfig() *oauth2.Config {
	return p.Config
}

func (p *GenericProvider) GetName() string {
	return p.Name
}
