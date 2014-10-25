# tt_goji_oauth
Add login with OAUTH to your GOJI webapp.

Assumes there's a redigo redis connection in c.Env["redis"]. Or could just use a global map

## TODO
- GET basic login page
- POST login action
  x choose which provider
  x add a cookie for a session
  x create random state string, store in session
  x redirect to auth URL
- POST/GET?? oauth callback
  x extract random string from session and check
  x exchange temp code for token
  x get user info (user id, email)
  - get / create user info
  - update session to include user info / logged in status

  - move session out into middleware. Session needs to exist outside of the oauth code
  - flesh out DB thingy

/login - plain page for logging in.
/login/oauth/start
/login/oauth/callback