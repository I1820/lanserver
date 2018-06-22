package main

import (
	"log"

	"github.com/aiotrc/lanserver_sh/actions"
)

func main() {
	app := actions.App()
	if err := app.Serve(); err != nil {
		log.Fatal(err)
	}
}
