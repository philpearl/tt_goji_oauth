package views

import (
	"fmt"
	"log"
	"net/http"

	// "github.com/golang/oauth2"
	"github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/providers"
	"github.com/zenazn/goji/web"
)

func StartLogin(c web.C, w http.ResponseWriter, r *http.Request) {

	sh := c.Env["sessionholder"].(base.SessionHolder)
	if sh == nil {
		log.Panicf(`misconfiguration - no session holder in c.Env["sessionholder"]`)
	}
	providerStore := c.Env["providerstore"].(*providers.ProviderStore)
	if providerStore == nil {
		log.Panicf(`misconfiguration - no provider store in c.Env["providerstore"]`)
	}

	// Get a session
	// TODO: what does it mean to login if we have a session already and it is marked as logged in?
	var s *base.Session
	si := c.Env["session"]
	if si == nil {
		s = sh.Create(c)
	} else {
		s = si.(*base.Session)
	}
	sh.AddToResponse(c, s, w)

	// next parameter holds url we're aiming for.  We store this in the session for later
	r.ParseForm()
	next := r.Form.Get("next")
	if next != "" {
		s.Put("next", next)
	}

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
	conf := provider.GetConfig()

	// Create our oauth state.  This includes a random secret we check later.  This is stored in the session
	state := newOauthState()
	state.ProviderName = providerName
	s.Put("oauth:secret", state.Secret)

	url := conf.AuthCodeURL(state.Encode(), "online", "auto")
	log.Printf("redirect to %s", url)

	h := w.Header()
	h.Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
