package base

import (
	"github.com/philpearl/tt_goji_middleware/base"
	"github.com/philpearl/tt_goji_oauth/providers"
)

/*
 Context for tt_goji_oauth
*/
type Context struct {
	SessionHolder base.SessionHolder
	ProviderStore *providers.ProviderStore
	Callbacks     Callbacks
}
