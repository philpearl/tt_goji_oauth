package tt_goji_oauth

import (
	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/philpearl/tt_goji_oauth/providers"
	"github.com/philpearl/tt_goji_oauth/views"

	"github.com/zenazn/goji/web"
)

/*
Build a mux to handle the endpoints required for oauth.

 baseUrl   - the full URL to the point where the oauth views are added,
             including trailing /
 prefix    - path prefix for the oauth views without trailing /
 sessionHolder - the session store used with th tt_goji_middleware session middleware
 callbacks - callbacks from tt_goji_oauth to application code
 providers - providers to include in the store.

For example "http://localhost:7778/login/oauth/", "/login/oauth".  As you can tell we only have two
parameters because I've been too lazy to parse the url

This function assumes that the session middleware from [tt_goji_middleware](https://github.com/philpearl/tt_goji_middleware)
is in the stack.
*/
func Build(baseUrl, prefix string, sessionHolder mbase.SessionHolder, callbacks base.Callbacks, provs ...func(baseUrl string) providers.Provider) *web.Mux {
	context := &base.Context{
		SessionHolder: sessionHolder,
		ProviderStore: providers.NewProviderStore(baseUrl, provs...),
		Callbacks:     callbacks,
	}

	mux := web.New()

	// Add a handler to strip the first part of the url so we don't need to
	// match it for all endpoints
	mux.Use(mbase.BuildStripPrefix(prefix))

	mux.Use(mbase.BuildEnvSet("oauth:context", context))

	mux.Post("/start/:provider/", views.StartLogin)

	mux.Get("/callback/", views.OauthCallback)
	mux.Post("/logout/", views.Logout)
	mux.Compile()

	return mux
}
