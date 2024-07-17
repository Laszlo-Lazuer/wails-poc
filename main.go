package main

import (
	"embed"

	"github.com/sanbornm/go-selfupdate/selfupdate"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// the app's version. This will be set on build.
var version string

// go-selfupdate setup and config
var updater = &selfupdate.Updater{
	CurrentVersion: version,                                                      // Manually update the const, or set it using `go build -ldflags="-X main.VERSION=<newver>" -o hello-updater src/hello-updater/main.go`
	ApiURL:         "https://laszlo-lazuer.github.io/artifactory/hello-updater/", // The server hosting `$CmdName/$GOOS-$ARCH.json` which contains the checksum for the binary
	BinURL:         "http://localhost:8080/",                                     // The server hosting the zip file containing the binary application which is a fallback for the patch method
	DiffURL:        "http://localhost:8080/",                                     // The server hosting the binary patch diff for incremental updates
	Dir:            "update/",                                                    // The directory created by the app when run which stores the cktime file
	CmdName:        "basic",                                                      // The app name which is appended to the ApiURL to look for an update
	ForceCheck:     true,                                                         // For this example, always check for an update unless the version is "dev"
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "basic",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
