package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx       context.Context
	htmlFiles [][]byte // a slice of byte slices
	currPage  int
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) LoadEpubFile() {
	fmt.Println("opening dialog options")

	dialogOptions := runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				Pattern: "*.epub",
			},
		},
	}

	filepath, err := runtime.OpenFileDialog(a.ctx, dialogOptions)
	if err != nil {
		log.Fatal(err)
	}

	if filepath == "" {
		return
	}

	a.currPage = 0

	a.htmlFiles = [][]byte{}

	// open epub file as zip file
	r, err := zip.OpenReader(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// TODO: are we sure this will capture only files with
	// extension .xhtml, .htm, .html?
	htmlRegex := regexp.MustCompile(`.x{0,1}html{0,1}`)
	cssRegex := regexp.MustCompile(`.css`) // and this too for .css?

	for _, f := range r.File {
		if htmlRegex.MatchString(f.Name) {
			fmt.Printf("found html file: %s\n", f.Name)

			fileReader, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			data, err := io.ReadAll(fileReader)
			if err != nil {
				log.Fatal(err)
			}

			a.htmlFiles = append(a.htmlFiles, data)

			fileReader.Close()
		} else if cssRegex.MatchString(f.Name) {
			fileReader, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			cssData, err := io.ReadAll(fileReader)
			if err != nil {
				log.Fatal(err)
			}

			runtime.EventsEmit(a.ctx, "style", string(cssData))

			fileReader.Close()
		}
	}

	// for now, just show the first html file
	runtime.EventsEmit(a.ctx, "page", string(a.htmlFiles[0]))
}

func (a *App) PrevFile() {
	if len(a.htmlFiles) == 0 {
		return
	}

	a.currPage = a.currPage - 1
	if a.currPage < 0 {
		a.currPage = 0
	}

	runtime.EventsEmit(a.ctx, "page", string(a.htmlFiles[a.currPage]))
}

func (a *App) NextFile() {
	if len(a.htmlFiles) == 0 {
		return
	}

	a.currPage = a.currPage + 1
	if a.currPage > len(a.htmlFiles)-1 {
		a.currPage = 0
	}

	runtime.EventsEmit(a.ctx, "page", string(a.htmlFiles[a.currPage]))
}
