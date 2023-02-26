package main

type App struct {
	Server string
}

func NewApp() *App {

	return &App{
		Server: "test",
	}
}
