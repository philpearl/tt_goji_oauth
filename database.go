package tt_goji_oauth

type Database interface {
	GetOrCreateUser(map[string]interface{})
}
