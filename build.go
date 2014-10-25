package tt_goji_oauth

import (
	"github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/providers"
	"github.com/philpearl/tt_goji_oauth/redis"
	"github.com/philpearl/tt_goji_oauth/views"

	"github.com/zenazn/goji/web"
)

/*
Build a mux to handle the endpoints required for oauth.

 baseUrl - the full URL to the point where the oauth views are added, including trailing /
 prefix  - path prefix for the oauth views without trailing /

For example "http://localhost:7778/login/oauth/", "/login/oauth".  As you can tell we only have two
parameters because I've been too lazy to parse the url

This function assumes that there are already middleware set up to add a redis connection
to c.Env["redis"], and that the session middleware from github.com/philpearl/tt_goji_middleware
is in the stack.
*/
func Build(baseUrl, prefix string) *web.Mux {

	sessionHolder := redis.NewSessionHolder()
	providerStore := providers.NewProviderStore(baseUrl, providers.Github)

	mux := web.New()

	// Add a handler to strip the first part of the url so we don't need to
	// match it for all endpoints
	mux.Use(base.BuildStripPrefix(prefix))

	mux.Use(base.BuildEnvSet("sessionholder", sessionHolder))
	mux.Use(base.BuildEnvSet("providerstore", providerStore))

	mux.Get("/start/:provider/", views.StartLogin) // TODO - should be Post - Get easier for immediate testing
	mux.Get("/callback/", views.OauthCallback)
	mux.Compile()

	return mux
}
