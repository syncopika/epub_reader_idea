package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"fmt"
	"net/http"
	"os"
	"strings"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	app *App
	http.Handler
}

func NewFileLoader(a *App) *FileLoader {
	f := &FileLoader{}
	f.app = a
	return f
}

// we should expect any assets for this app (e.g. the images for the e-book)
// to be in a specific temporary directory specified via app.tmpDirectoryName
func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	requestedFilenameParts := strings.Split(requestedFilename, "/")
	requestedFilename = requestedFilenameParts[len(requestedFilenameParts)-1]

	fileData, err := os.ReadFile(fmt.Sprintf("%s/%s", h.app.tmpDirectoryName, requestedFilename))
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "epub_reader_idea",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(app),
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
