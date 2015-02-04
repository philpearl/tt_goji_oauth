/*
View functions for OAUTH
*/
package views

import (
	"fmt"
	"log"
	"net/http"

	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

/*
StartLogin kicks off the OAUTH process.

Add a 'next' url parameter to control where the user is redirected to after a successful login.
*/
func StartLogin(c web.C, w http.ResponseWriter, r *http.Request) {
	context := c.Env["oauth:context"].(*base.Context)
	sh := context.SessionHolder
	providerStore := context.ProviderStore

	// Get a session
	// TODO: what does it mean to login if we have a session already and it is marked as logged in?
	var s *mbase.Session
	si := c.Env["session"]
	if si == nil {
		s = sh.Create(c)
	} else {
		s = si.(*mbase.Session)
	}
	sh.AddToResponse(c, s, w)

	// next parameter holds url we're aiming for.  We store this in the session for later
	r.ParseForm()
	next := r.Form.Get("next")
	s.Put("next", next)

	// Redirect the user to the appropriate provider url
	providerName, ok := c.URLParams["provider"]
	if !ok {
		http.Error(w, "URL must contain an OAUTH provider name", http.StatusBadRequest)
		return
	}
	provider, ok := providerStore.GetProvider(providerName)
	if !ok {
		http.Error(w, fmt.Sprintf("Requested OAUTH provider %s not available", providerName), http.StatusNotFound)
		return
	}
	conf := provider.GetConfig(r)

	// Create our oauth state.  This includes a random secret we check later.  This is stored in the session
	state := newOauthState()
	state.ProviderName = providerName
	s.Put("oauth:secret", state.Secret)

	url := conf.AuthCodeURL(state.encode(), oauth2.AccessTypeOffline)
	log.Printf("redirect to %s", url)

	h := w.Header()
	h.Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
