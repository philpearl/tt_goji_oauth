# tt_goji_oauth
Add login with OAUTH to your [Goji](https://github/zenazn/goji) webapp.

[![Build Status](https://travis-ci.org/philpearl/tt_goji_oauth.svg)](https://travis-ci.org/philpearl/tt_goji_oauth) [![GoDoc](https://godoc.org/github.com/philpearl/tt_goji_oauth?status.svg)](https://godoc.org/github.com/philpearl/tt_goji_oauth)


Currently supports github, but additional OAUTH providers can be plugged in.

## Code state
I've only just written this and haven't used it in anger yet.

## How to use
See the example code in /example for full details, but the basics are as follows.

1. Create a SessionHolder from https://github.com/philpearl/tt_goji_middleware/base, and add session middleware to your mux.
2. Call tt_goji_oauth.Build() and add the handler it returns to your mux.  We suggest you add it at /login/oauth.
3. Add pages that have logged-in users to your mux beneath the session middleware.
4. To login, POST to /login/oauth/start/github/.  Add a 'next' parameter to control where the user is redirected to after login
5. When login completes the callback you registered calling tt_goji_oauth.Build() will be called with user information.  You should check the user against your database at this point, and set up information in the session.
6. The user will be redirected to /, or where-ever you specified via the next parameter.

## Future

1. More providers - including google & facebook that do client-side oauth flows.
2. Redirect to a different page for a new user.
3. User specified providers