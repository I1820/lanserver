package grifts

import (
	"github.com/aiotrc/lanserver/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
