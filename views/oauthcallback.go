package views

import (
	"log"
	"net/http"

	// "github.com/golang/oauth2"
	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/zenazn/goji/web"
)

func OauthCallback(c web.C, w http.ResponseWriter, r *http.Request) {
	context := c.Env["oauth:context"].(*base.Context)

	// Get the session
	var session *mbase.Session
	s, ok := c.Env["session"]
	log.Printf("env is %v", c.Env)
	if ok {
		log.Printf("session in env")
		session, ok = s.(*mbase.Session)
	}

	if !ok {
		log.Printf("Could not retrieve session in oauth callback.")
		http.Error(w, "no session", http.StatusBadRequest)
		return
	}

	// Important to check we've been passed back our random value
	val, ok := session.Get("oauth:secret")
	if !ok {
		log.Printf("No secret available in oauth callback")
		http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
		return
	}
	secret := val.(int64)

	// Extract interesting parameters from response.
	r.ParseForm()
	state, err := newOauthStateFromString(r.Form.Get("state"))
	if err != nil {
		log.Printf("Failed to decode returned oauth state. %v", err)
		http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
		return
	}

	if state.Secret != secret {
		log.Printf("Mismatched state in oauth callback.  Have %s expected %s", secret, state.Secret)
		http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
		return
	}

	// Get an access token in exchange for our temporary token
	// In fact we get a "transport" that can issue requests with authentication, then
	// ask for some user info.  This will get our token as a side effect
	// Not that we really need a token - we just want to identify our user
	rcode := r.Form.Get("code")
	provider, ok := context.ProviderStore.GetProvider(state.ProviderName)
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
	err = context.Callbacks.GetOrCreateUser(c, state.ProviderName, user)
	if err != nil {
		log.Printf("GetOrCreateUser callback failed. %v", err)
		http.Error(w, "", http.StatusServiceUnavailable)
		return
	}

	// Mark the session as logged in
	session.Put("logged_in", true)

	// Redirect to final destination.
	val, ok = session.Get("next")
	var url string
	if !ok || val == "" {
		url = "/"
	} else {
		url = val.(string)
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
