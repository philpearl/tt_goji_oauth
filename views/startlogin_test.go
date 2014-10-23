package views

import (
	"net/http"
	"net/http/httptest"
	"testing"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji/web"

	"github.com/philpearl/tt_goji_oauth/providers"
	"github.com/philpearl/tt_goji_oauth/redis"
)

func TestStartLogin(t *testing.T) {
	// Need a providerstore in c.Env["providerstore"] and a sessionHolder in c.Env["sessionholder"]
	sessionHolder := redis.NewSessionHolder()
	providerStore := providers.NewProviderStore("http://localhost/")

	providerStore.Add(
		"github",
		"dummy",
		"dummy",
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token",
		[]string{"user"},
	)

	conn, err := redigo.Dial("tcp", ":6379")
	if err != nil {
		t.Skipf("Cannot connect to redis. %v", err)
	}
	c := web.C{
		Env: map[string]interface{}{
			"redis":         conn,
			"sessionholder": sessionHolder,
			"providerstore": providerStore,
		},
	}

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	StartLogin(c, w, r)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302? got %d", w.Code)
	}

	url := w.HeaderMap.Get("Location")
	s := sessionHolder.Get(c, r)

	if url == "" {
		t.Fatalf("expected redirect to callback, url is empty")
	}
}
