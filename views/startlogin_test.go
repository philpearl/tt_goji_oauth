package views

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji/web"

	mbase "github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_middleware/redis"
	"github.com/philpearl/tt_goji_oauth/base"
	"github.com/philpearl/tt_goji_oauth/providers"
)

func buildTestContext(t *testing.T) *web.C {
	// Need a providerstore in c.Env["providerstore"] and a sessionHolder in c.Env["sessionholder"]
	sessionHolder := redis.NewSessionHolder()
	os.Setenv("GITHUB_CLIENT_ID", "dummy")
	os.Setenv("GITHUB_CLIENT_SECRET", "dummy_secret")
	providerStore := providers.NewProviderStore("http://localhost/", providers.Github)

	context := &base.Context{
		SessionHolder: sessionHolder,
		ProviderStore: providerStore,
	}

	conn, err := redigo.Dial("tcp", ":6379")
	if err != nil {
		t.Skipf("Cannot connect to redis. %v", err)
	}
	c := web.C{
		Env: map[interface{}]interface{}{
			"redis":         conn,
			"oauth:context": context,
		},
		URLParams: map[string]string{
			"provider": "github",
		},
	}

	return &c
}

func TestStartLogin(t *testing.T) {
	c := buildTestContext(t)

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	StartLogin(*c, w, r)

	if w.Code != http.StatusFound {
		t.Fatalf("expected 302? got %d, %s", w.Code, w.Body.String())
	}

	url := w.HeaderMap.Get("Location")

	if url == "" {
		t.Fatalf("expected redirect to callback, url is empty")
	}

	s := c.Env["session"].(*mbase.Session)
	secret, ok := s.Get("oauth:secret")
	if !ok {
		t.Fatalf("failed to save secret %v", secret)
	}
}

func TestStartLoginNoProvider(t *testing.T) {
	c := buildTestContext(t)
	delete(c.URLParams, "provider")

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	StartLogin(*c, w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d, %s", w.Code, w.Body.String())
	}
}

func TestStartLoginBadProvider(t *testing.T) {
	c := buildTestContext(t)
	c.URLParams["provider"] = "cheese"

	r, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()
	StartLogin(*c, w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 400 got %d, %s", w.Code, w.Body.String())
	}
}

func TestStartLoginNext(t *testing.T) {
	c := buildTestContext(t)

	v := url.Values{}
	v.Set("next", "http://example.com/next")

	r, _ := http.NewRequest("POST", "/?"+v.Encode(), nil)
	w := httptest.NewRecorder()
	StartLogin(*c, w, r)

	if w.Code != http.StatusFound {
		t.Fatalf("expected redirect got %d, %s", w.Code, w.Body.String())
	}

	s := c.Env["session"].(*mbase.Session)

	next, ok := s.Get("next")
	if !ok || next.(string) != "http://example.com/next" {
		t.Fatalf("next url not stored. %v, %v", ok, next)
	}
}
