package grifts

import (
	"github.com/aiotrc/lanserver_sh/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
