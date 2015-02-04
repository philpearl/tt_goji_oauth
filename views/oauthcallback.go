package views

import (
	"log"
	"net/http"

	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
)

/*
When the OAUTH service is done it redirects the user to this view with
a temporary code that it can exchange for an OAUTH token.  It then calls
the provider to get identifying information for the user.
*/
func OauthCallback(c web.C, w http.ResponseWriter, r *http.Request) {
	context := c.Env["oauth:context"].(*base.Context)

	// Get the session
	session, ok := mbase.SessionFromEnv(&c)

	if !ok {
		log.Printf("Could not retrieve session in oauth callback.")
		http.Error(w, "no session", http.StatusBadRequest)
		return
	}

	// Extract interesting parameters from response.
	r.ParseForm()
	haveState := false
	providerName := r.Form.Get("provider")
	if providerName == "" {
		stateParam := r.Form.Get("state")
		if stateParam == "" {
			log.Printf("State parameter not included")
			http.Error(w, "OAUTH state parameter not present", http.StatusBadRequest)
			return
		}
		state, err := newOauthStateFromString(stateParam)
		if err != nil {
			log.Printf("Failed to decode returned oauth state. %v", err)
			http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
			return
		}

		// Important to check we've been passed back our random value.
		val, ok := session.Get("oauth:secret")
		if !ok {
			log.Printf("No secret available in oauth callback")
			http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
			return
		}
		secret := val.(int64)

		if state.Secret != secret {
			log.Printf("Mismatched state in oauth callback.  Have %s expected %s", secret, state.Secret)
			http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
			return
		}
		providerName = state.ProviderName
		haveState = true
	}

	// Get an access token in exchange for our temporary token
	// In fact we get a "transport" that can issue requests with authentication, then
	// ask for some user info.  This will get our token as a side effect
	// Not that we really need a token - we just want to identify our user

	provider, ok := context.ProviderStore.GetProvider(providerName)
	if !ok {
		log.Printf("Provider %s not found", providerName)
		http.Error(w, "Oops!", http.StatusInternalServerError)
		return
	}

	if provider.NeedState() && !haveState {
		// We expected the state parameter, but it did not arrive.  Security issue
		log.Printf("Expected state parameter not present")
		http.Error(w, "OAUTH protocol error detected", http.StatusBadRequest)
		return
	}

	conf := provider.GetConfig(r)
	rcode := r.Form.Get("code")
	token, err := conf.Exchange(oauth2.NoContext, rcode)
	if err != nil {
		log.Printf("Could not exchange code %s for a token. %v", rcode, err)
		http.Error(w, "Authentication error", http.StatusBadGateway)
		return
	}

	// Get some user info.
	client := conf.Client(oauth2.NoContext, token)
	user, err := provider.GetUserInfo(r, client, token)
	if err != nil {
		log.Printf("Failed to get information from user.  %v", err)
		http.Error(w, "Couldn't retrieve user info", http.StatusServiceUnavailable)
		return
	}
	log.Printf("Have user info %v", user)

	// Get or Create a user object.  Again, some kind of plug-in storage would make sense
	url, err := context.Callbacks.GetOrCreateUser(c, providerName, user)
	if err != nil {
		log.Printf("GetOrCreateUser callback failed. %v", err)
		http.Error(w, "", http.StatusServiceUnavailable)
		return
	}

	// Mark the session as logged in
	session.Put("logged_in", true)

	// Redirect to final destination.  We use the value returned by GetOrCreateUser if set
	if url == "" {
		val, ok := session.Get("next")
		if !ok || val == "" {
			url = "/"
		} else {
			url = val.(string)
			// Don't reuse next
			session.Del("next")
		}
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
