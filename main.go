// D75 Studio: configurador del teclado Weikav D75 para Linux y Windows.
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "D75 Studio",
		Description: "Configurador del teclado Weikav D75",
		Services: []application.Service{
			application.NewService(&KeyboardService{}),
			application.NewService(&ProfileService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "D75 Studio",
		Width:            1380,
		Height:           860,
		MinWidth:         1180,
		MinHeight:        760,
		BackgroundColour: application.NewRGB(11, 12, 16),
		URL:              "/",
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
