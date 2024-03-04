package bootstrap

import (
	"fmt"
	"log"
)

func (app *App) Start() {
	fmt.Printf("[socket-processor] service started...\n\n\n\n")
	err := app.router.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
