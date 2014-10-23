package views

import (
	"log"
	"net/http"

	// "github.com/golang/oauth2"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/philpearl/tt_goji_oauth/providers"
	"github.com/zenazn/goji/web"
)

func OauthCallback(c web.C, w http.ResponseWriter, r *http.Request) {
	sh := c.Env["sessionholder"].(base.SessionHolder)
	providerStore := c.Env["providerstore"].(*providers.ProviderStore)

	// Get the session
	s, err := sh.Get(c, r)
	if err != nil {
		log.Printf("Could not retrieve session in oauth callback. %v", err)
		http.Error(w, "no session", http.StatusBadRequest)
		return
	}

	// Important to check we've been passed back our random value
	val, ok := s.Get("oauth:state")
	if !ok {
		log.Printf("No state available in oauth callback")
		http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
		return
	}
	state := val.(string)

	// Extract interesting parameters from response.
	r.ParseForm()
	rstate := r.Form.Get("state")
	if rstate != state {
		log.Printf("Mismatched state in oauth callback.  Have %s expected %s", rstate, state)
		http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
		return
	}

	// Get an access token in exchange for our temporary token
	// In fact we get a "transport" that can issue requests with authentication, then
	// ask for some user info.  This will get our token as a side effect
	// Not that we really need a token - we just want to identify our user
	rcode := r.Form.Get("code")
	provider, ok := providerStore.GetProvider("github")
	if !ok {
		http.Error(w, "Oops!", http.StatusInternalServerError)
		return
	}
	conf := provider.GetConfig()
	token, err := conf.Exchange(rcode)
	if err != nil {
		http.Error(w, "Authentication error", http.StatusBadGateway)
		return
	}

	log.Printf("Here's our token %v", token)

	// Get some user info.  Note we're github only at this point - if we want to support other
	// providers here's where we would need to plug something in.
	t := conf.NewTransport()
	t.SetToken(token)
	user, err := provider.GetUserInfo(t)
	if err != nil {
		log.Printf("Failed to get information from user.  %v", err)
		http.Error(w, "Couldn't retrieve user info", http.StatusServiceUnavailable)
		return
	}
	log.Printf("Have user info %v", user)

	// Get or Create a user object.  Again, some kind of plug-in storage would make sense

	// Mark the session as logged in
	s.Put("logged_in", true)
	sh.Save(c, s)

	// Redirect to final destination.
	val, ok = s.Get("next")
	var url string
	if !ok {
		url = "/"
	} else {
		url = val.(string)
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
