package redis

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/philpearl/tt_goji_oauth/base"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/zenazn/goji/web"
)

type SessionHolder struct {
}

func NewSessionHolder() base.SessionHolder {
	return &SessionHolder{}
}

/*
   Get the session for this request
*/
func (sh *SessionHolder) Get(c web.C, r *http.Request) (base.Session, error) {
	cookie, err := r.Cookie("sessionid")
	if err != nil {
		return nil, err
	}
	log.Printf("Cookie is %v", cookie)
	sessionId := cookie.Value

	conn := c.Env["redis"].(redigo.Conn)

	sessionBytes, err := redigo.Bytes(conn.Do("GET", sessionKey(sessionId)))
	if err != nil {
		log.Printf("Could not find session for id %s", sessionId)
		return nil, err
	}

	dec := gob.NewDecoder(bytes.NewReader(sessionBytes))
	var session session
	err = dec.Decode(&session)

	return &session, err
}

/*
   Create a new session
*/
func (sh *SessionHolder) Create() base.Session {
	return &session{
		Idd:    generateSessionId(),
		Values: make(map[string]interface{}, 0),
	}
}

func generateSessionId() string {
	a := uint64(rand.Int63())
	b := uint64(rand.Int63())

	return strconv.FormatUint(a, 36) + strconv.FormatUint(b, 36)
}

/*
   Destroy a session so that it can no longer be retrieved
*/
func (sh *SessionHolder) Destroy(c web.C, session base.Session) error {
	sessionId := session.Id()
	conn := c.Env["redis"].(redigo.Conn)

	_, err := conn.Do("DEL", sessionKey(sessionId))
	return err
}

/*
   Save a session
*/
func (sh *SessionHolder) Save(c web.C, session base.Session) error {
	sessionId := session.Id()
	conn := c.Env["redis"].(redigo.Conn)

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	err := enc.Encode(session)
	if err != nil {
		return err
	}

	log.Printf("save session %s", sessionId)
	_, err = conn.Do("SET", sessionKey(sessionId), b.String())

	return err
}

type session struct {
	Idd    string
	Values map[string]interface{}
}

func (s *session) Id() string {
	return s.Idd
}

/*
   Add the session to the response

   Basically this means setting a cookie
*/
func (s *session) AddToResponse(w http.ResponseWriter) {
	cookie := http.Cookie{
		Name:   "sessionid",
		Value:  s.Idd,
		Path:   "/",
		MaxAge: 30 * 24 * 60 * 60,
	}
	http.SetCookie(w, &cookie)
}

/*
   Retrieve a value from the session
*/
func (s *session) Get(key string) (interface{}, bool) {
	val, ok := s.Values[key]
	return val, ok
}

/*
   Save a value in the session
*/
func (s *session) Put(key string, value interface{}) {
	s.Values[key] = value
}

func sessionKey(sessionId string) string {
	return fmt.Sprintf("sess:%s", sessionId)
}
