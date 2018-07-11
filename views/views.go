package views

import "github.com/philpearl/tt_goji_oauth/base"

type Views struct {
	cxt *base.Context
}

func New(cxt *base.Context) *Views {
	return &Views{
		cxt: cxt,
	}
}
