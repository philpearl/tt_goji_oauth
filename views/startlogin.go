package views

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"

	// "github.com/golang/oauth2"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/philpearl/tt_goji_oauth/providers"
	"github.com/zenazn/goji/web"
)

func StartLogin(c web.C, w http.ResponseWriter, r *http.Request) {

	sh := c.Env["sessionholder"].(base.SessionHolder)
	providerStore := c.Env["providerstore"].(*providers.ProviderStore)

	// Get a session
	s, err := sh.Get(c, r)
	if err != nil {
		s = sh.Create()
	}
	s.AddToResponse(w)

	// next parameter holds url we're aiming for
	r.ParseForm()
	next := r.Form.Get("next")
	if next != "" {
		log.Printf("adding next: %s", next)
		s.Put("next", next)
	}

	// Create the random state - we check this later so we save it in the session
	state := strconv.FormatUint(uint64(rand.Int63()), 36)
	s.Put("oauth:state", state)
	sh.Save(c, s)

	// Redirect the user to the appropriate provider url
	provider, ok := providerStore.GetProvider("github")
	if !ok {
		log.Printf("Failed to get OAUTH config for github")
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
	conf := provider.GetConfig()
	url := conf.AuthCodeURL(state, "online", "auto")
	log.Printf("redirect to %s", url)

	h := w.Header()
	h.Set("Location", url)
	w.WriteHeader(http.StatusFound)
}
