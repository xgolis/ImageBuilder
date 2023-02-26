package main

import (
	"ImageBuilder/cmd/ImageBuilder/app"
	"fmt"
)

func main() {
	app := app.NewApp()
	fmt.Printf(app.Server)
}
