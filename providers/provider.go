package providers

import (
	"net/http"

	"github.com/philpearl/oauth2"
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
	// GetName returns the name of the provider
	GetName() string
	// GetUserInfo gets information about the user.
	//
	// Information is returned in a map. Keys are defined by the PROVIDER_ constants
	GetUserInfo(r *http.Request, client *http.Client, t *oauth2.Token) (map[string]interface{}, error)

	// GetConfig returns the oauth2 config for this provider
	GetConfig(r *http.Request) *oauth2.Config

	NeedState() bool
}

/*
GenericProvider is a base partial implementation of Provider to use to build
full provider implementations.

To create a new provider:

1. Create a new provider struct containing GenericProvider

2. Implement GetUserInfo()

3. Create a function to create the provider and the embedded oauth config.

 type MyServiceProvider struct {
	GenericProvider
 }

 func (p *MyServiceProvider) GetUserInfo(t *oauth2.Transport) (map[string]interface{}, error) {
 	// Use t to make authenticated requests to My service to get user information and
 	// return it in a map
 }

 func MyService(baseUrl string) Provider {
	options := &oauth2.Options{
		ClientID:     ???,
		ClientSecret: ???,
		RedirectURL:  baseUrl + "callback/",
		Scopes:       []string{"user"},
	}
	config, _ := oauth2.NewConfig(options, authUrl, tokenUrl)

	return &MyServiceProvider{
		GenericProvider: GenericProvider{
			Name:   "MyService",
			Config: config,
		},
	}
 }

Ahh, except I don't currently have a great mechanism to let you add new providers
*/
type GenericProvider struct {
	Name   string
	Config *oauth2.Config
}

func (p *GenericProvider) GetConfig(r *http.Request) *oauth2.Config { return p.Config }
func (p *GenericProvider) GetName() string                          { return p.Name }
func (p *GenericProvider) NeedState() bool                          { return true }
