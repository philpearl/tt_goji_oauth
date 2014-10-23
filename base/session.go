package base

import (
	"github.com/zenazn/goji/web"
	"net/http"
)

type Session interface {
	/*
	   Get the id for this session
	*/
	Id() string

	/*
	   Add the session to the response
	*/
	AddToResponse(w http.ResponseWriter)

	/*
	   Retrieve a value from the session
	*/
	Get(key string) (interface{}, bool)

	/*
	   Save a value in the session
	*/
	Put(key string, value interface{})
}

type SessionHolder interface {
	/*
	   Get the session associated with the current request, if there is one.
	*/
	Get(c web.C, r *http.Request) (Session, error)

	/*
	   Create a new session
	*/
	Create() Session

	/*
	   Destroy a session so that it can no longer be retrieved
	*/
	Destroy(c web.C, session Session) error

	/*
	   Save a session
	*/
	Save(c web.C, session Session) error
}
