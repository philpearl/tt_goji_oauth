package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_middleware/redis"
	"github.com/philpearl/tt_goji_oauth"
)

var (
	// The template for /
	index *template.Template
)

func init() {
	index = template.Must(template.ParseFiles("index.html"))
}

/*
Callbacks implements the tt_goji_oauth callbacks.
*/
type Callbacks struct {
}

/*
GetOrCreateUser is called by tt_goji_oauth
*/
func (cbk Callbacks) GetOrCreateUser(c web.C, providerName string, user map[string]interface{}) (string, error) {
	// Here is where we should ensure the user info is stored in the DB, but we can cheat somewhat by
	// just adding the user info to the session
	session, _ := base.SessionFromEnv(&c)

	session.Put("user", user)
	return "", nil
}

/*
IndexView is the handler for /
*/
func IndexView(c web.C, w http.ResponseWriter, r *http.Request) {
	var user map[string]interface{}
	session, ok := base.SessionFromEnv(&c)
	if ok {
		// Our user is stored in the session if we're logged in
		var usr interface{}
		usr, ok = session.Get("user")
		if ok {
			log.Printf("Have user")
			user = usr.(map[string]interface{})
		} else {
			ok = false
		}
	}

	// Render the template
	index.Execute(w, struct {
		LoggedIn bool
		User     map[string]interface{}
	}{
		LoggedIn: ok,
		User:     user,
	})
}

func main() {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("caught a panic: %v\n", r)
			os.Exit(1)
		}
	}()

	if !flag.Parsed() {
		flag.Parse()
	}

	// Build a mux for the site.
	m := web.New()
	m.Use(base.LoggingMiddleWare)
	m.Use(middleware.Recoverer)
	m.Use(middleware.EnvInit)

	// tt_goji_oauth requires redis and sessions middleware.
	m.Use(redis.BuildRedis(":6379"))
	sh := redis.NewSessionHolder()
	m.Use(base.BuildSessionMiddleware(sh))

	// Add the oauth mux under /login/oauth/
	callbacks := Callbacks{}
	oauthm, _ := tt_goji_oauth.Build("http://localhost:8000/login/oauth/", "/login/oauth", sh, callbacks)
	m.Handle("/login/oauth/*", oauthm)

	// Add our main page in the root
	m.Get("/", IndexView)
	m.Compile()

	// Now add goji boilerplate to run the site
	http.Handle("/", m)
	listener := bind.Default()
	log.Println("Starting Goji on", listener.Addr())

	bind.Ready()

	err := graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		log.Fatal(err)
	}

	graceful.Wait()
}
