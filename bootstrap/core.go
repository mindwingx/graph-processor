package bootstrap

import "github.com/mindwingx/graph-processor/driver"

func (app *App) initRegistry() {
	app.registry = driver.NewViper()
	app.registry.InitRegistry()
}

func (app *App) initRouter() {
	app.router = driver.InitRouter(app.registry)
	app.router.InitWsConnWorkerPool()
	app.router.Routes()
}
