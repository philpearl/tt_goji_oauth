package views

import (
	"net/http"

	// "github.com/golang/oauth2"
	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/zenazn/goji/web"
)

/*
Logout deletes the current session

Add a 'next' url parameter to control where the user is redirected to after logout.
*/
func Logout(c web.C, w http.ResponseWriter, r *http.Request) {
	// Get the session
	session, ok := mbase.SessionFromEnv(&c)

	if ok {
		context := c.Env["oauth:context"].(*base.Context)
		sh := context.SessionHolder
		sh.Destroy(c, session)
	}

	r.ParseForm()
	url := r.Form.Get("next")
	if url != "" {
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
	}
}
