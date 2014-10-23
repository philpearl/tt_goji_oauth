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

This assumes that there are already middleware set up to add a redis connection
to c.Env
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

	mux.Get("/start/", views.StartLogin) // TODO - should be Post
	mux.Get("/callback/", views.OauthCallback)
	mux.Compile()

	return mux
}
