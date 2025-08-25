package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/ironpark/teatime/internal/database"
	"github.com/ironpark/teatime/services"
	"github.com/ironpark/teatime/stores"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	db, err := database.OpenInMemory()
	if err != nil {
		log.Fatal(err)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	recipesDir := filepath.Join(home, ".teatime", "recipes")
	if err := os.MkdirAll(recipesDir, 0755); err != nil {
		log.Fatal(err)
	}
	store := stores.NewStore(db, recipesDir)

	// Initialize services
	recipesService := services.NewRecipesService(store)
	triggersService := services.NewTriggersService(store, recipesService)

	app := application.New(application.Options{
		Name:        "teatime",
		Description: "Teatime is a local workflow engine",
		Services: []application.Service{
			application.NewService(recipesService),
			application.NewService(triggersService),
			application.NewService(services.NewSettingsService(store)),
			application.NewService(services.NewSecretsService(store)),
			application.NewService(services.NewEnvironmentVariablesService(store)),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "TeaTime",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
	// Start triggers service
	if err := triggersService.Start(app.Context()); err != nil {
		log.Printf("Warning: Failed to start triggers service: %v", err)
	}
	err = app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
