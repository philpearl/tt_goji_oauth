package base

import (
	"github.com/zenazn/goji/web"
)

/*
Callbacks defines the interface between tt_goji_oauth and the ret of the application
*/
type Callbacks interface {
	/*
	   GetOrCreateUser is called when a user logs in and tt_goji_oauth has obtained
	   identifying information about the user from the oauth provider.
	*/
	GetOrCreateUser(c web.C, providerName string, user map[string]interface{}) error
}
