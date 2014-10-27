package views

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"log"
	"math/rand"
	"strings"
)

type oauthState struct {
	Secret       int64
	ProviderName string
}

func newOauthState() *oauthState {
	return &oauthState{
		Secret: rand.Int63(),
	}
}

func newOauthStateFromString(encoded string) (*oauthState, error) {
	r := strings.NewReader(encoded)
	b64Dec := base64.NewDecoder(base64.URLEncoding, r)
	dec := gob.NewDecoder(b64Dec)
	var state oauthState
	err := dec.Decode(&state)

	return &state, err
}

func (o *oauthState) encode() string {
	var b bytes.Buffer
	b64Enc := base64.NewEncoder(base64.URLEncoding, &b)
	enc := gob.NewEncoder(b64Enc)
	err := enc.Encode(o)
	if err != nil {
		// Should be no reason for this to happen outside of dev
		log.Panicf("Could not encode oauth state.  %v", err)
	}
	// Must close the base64 encode to flush out everything
	b64Enc.Close()

	return b.String()
}
