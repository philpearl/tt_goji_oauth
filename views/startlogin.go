/*
View functions for OAUTH
*/
package views

import (
	"fmt"
	"log"
	"net/http"

	"github.com/philpearl/oauth2"
	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/zenazn/goji/web"
)

/*
StartLogin kicks off the OAUTH process.

Add a 'next' url parameter to control where the user is redirected to after a successful login.
*/
func (v *Views) StartLogin(c web.C, w http.ResponseWriter, r *http.Request) {

	// Redirect the user to the appropriate provider url
	providerName, ok := c.URLParams["provider"]
	if !ok {
		http.Error(w, "URL must contain an OAUTH provider name", http.StatusBadRequest)
		return
	}

	url := v.StartLoginURL(c, w, r, providerName)
	log.Printf("redirect to %s", url)

	h := w.Header()
	h.Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

// StartLoginURL returns the URL to redirect to to start the OAUTH process. It also alters the
// session to hold login info, and updates cookies on the response. This is used when you need
// to redirect to the OAUTH URL via "unconventional" means - as for Shopify embedded apps
func (v *Views) StartLoginURL(c web.C, w http.ResponseWriter, r *http.Request, providerName string) string {
	sh := v.cxt.SessionHolder
	providerStore := v.cxt.ProviderStore

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

	provider, ok := providerStore.GetProvider(providerName)
	if !ok {
		http.Error(w, fmt.Sprintf("Requested OAUTH provider %s not available", providerName), http.StatusNotFound)
		return ""
	}
	conf := provider.GetConfig(r)

	// Create our oauth state.  This includes a random secret we check later.  This is stored in the session
	state := newOauthState()
	state.ProviderName = providerName
	s.Put("oauth:secret", state.Secret)

	return conf.AuthCodeURL(state.encode(), oauth2.AccessTypeOffline)
}
