package di

import "Week04/internal/server/http"

type App struct {
	httpServer *http.HttpServer
}

func NewApp(httpServer *http.HttpServer) (app *App, cf func(), err error) {
	app = &App{httpServer: httpServer}

	cf = func() {
		httpServer.Close()
	}
	return
}
