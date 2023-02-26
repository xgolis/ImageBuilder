package app

import (
	"net/http"
	"time"

	"github.com/xgolis/ImageBuilder/builder"
)

type App struct {
	Server *http.Server
}

func NewApp() *App {
	builder.Testik()

	mux := MakeHandlers()

	return &App{
		Server: &http.Server{
			Addr:           ":8080",
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}
