package redis

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji/web"
)

func TestSessionCreate(t *testing.T) {
	conn, err := redigo.Dial("tcp", ":6379")
	if err != nil {
		t.Skipf("Cannot connect to redis. %v", err)
	}
	c := web.C{
		Env: map[string]interface{}{"redis": conn},
	}
	sh := NewSessionHolder()

	// Create the session
	s := sh.Create()

	// Check we can read and write to the session
	s.Put("cheese", "cheddar")
	val, ok := s.Get("cheese")
	cheese := val.(string)

	if !ok || cheese != "cheddar" {
		t.Fatalf("expected cheddar, but %v, %v", ok, cheese)
	}

	_, ok = s.Get("hat")
	if ok {
		t.Fatalf("expected to not be OK")
	}

	err = sh.Save(c, s)
	if err != nil {
		t.Fatalf("failed to save session - %v", err)
	}

	// Check we can write the session Id to the response
	w := httptest.NewRecorder()
	s.AddToResponse(w)

	ch := w.HeaderMap.Get("Set-Cookie")
	if ch != fmt.Sprintf("sessionid=%s; Max-Age=2592000", s.Id()) {
		t.Fatalf("cookie header is %v", ch)
	}

	// Can we extract the session from a request carrying a cookie?
	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  "sessionid",
		Value: s.Id(),
	})

	s1, err := sh.Get(c, r)
	if err != nil {
		t.Fatalf("got error reading session from request, %v", err)
	}

	val, ok = s1.Get("cheese")
	cheese = val.(string)
	if !ok || cheese != "cheddar" {
		t.Fatalf("could not get right cheese %v %v", ok, cheese)
	}

	sh.Destroy(c, s1)

	s1, err = sh.Get(r)
	if err == nil {
		t.Fatalf("session should be destroyed")
	}

}
