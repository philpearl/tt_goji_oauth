# tt_goji_oauth
Add login with OAUTH to your GOJI webapp.

Assumes there's a redigo redis connection in c.Env["redis"]. Or could just use a global map

## TODO
- GET basic login page
- POST login action
  - choose which provider
  - add a cookie for a session
  - create random "" string, store in session
  - redirect to auth URL
- POST/GET?? oauth callback
  - extract random string from session and check
  - exchange temp code for token
  - get user info (user id, email)
  - get / create user info
  - update session to include user info / logged in status