package app

import "github.com/xgolis/ImageBuilder/builder"

type App struct {
	Server string
}

func NewApp() *App {
	builder.Testik()

	return &App{
		Server: "test",
	}
}
