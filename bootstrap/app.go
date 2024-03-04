package bootstrap

import "github.com/mindwingx/graph-processor/driver/abstractions"

type App struct {
	registry abstractions.RegAbstraction
	router   abstractions.RouterAbstraction
}

func NewApp() *App {
	return &App{}
}

func (app *App) Init() {
	app.initRegistry()
	app.initRouter()
}
