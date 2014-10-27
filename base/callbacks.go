/*
Package base contains interface and struct definitions used by the rest of the code.
*/
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

	   providerName - Name of the provider that produced the user information (e.g. "github")
	   user         - Map of user information from the provider. see github.com/philpearl/tt_goji_oauth/providers

	   Return a url if you want to redirect the user to a particular page (e.g. to fill in a user profile form
	   for a new user, or to read updated terms & conditions)
	*/
	GetOrCreateUser(c web.C, providerName string, user map[string]interface{}) (string, error)
}
