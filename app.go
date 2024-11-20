package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx              context.Context
	htmlFiles        [][]byte
	currPage         int
	tmpDirectoryName string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.tmpDirectoryName = "epub_reader_tmp_dir"
	os.Mkdir(a.tmpDirectoryName, os.ModePerm)
}

// on shutdown
func (a *App) shutdown(ctx context.Context) {
	os.RemoveAll(a.tmpDirectoryName)
}

func (a *App) HtmlRegex() *regexp.Regexp {
	return regexp.MustCompile(`\.x{0,1}html{0,1}$`)
}

func (a *App) CssRegex() *regexp.Regexp {
	return regexp.MustCompile(`\.css$`)
}

func (a *App) ImageRegex() *regexp.Regexp {
	return regexp.MustCompile(`(\.png|.jpg)$`)
}

func (a *App) LoadEpubFile() {
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

	htmlRegex := a.HtmlRegex()
	cssRegex := a.CssRegex()
	imageRegex := a.ImageRegex()

	for _, f := range r.File {
		if htmlRegex.MatchString(f.Name) {
			//fmt.Printf("found html file: %s\n", f.Name)

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
		} else if imageRegex.MatchString(f.Name) {
			// write image to temp directory
			fileReader, err := f.Open()

			if err != nil {
				log.Fatal(err)
			}

			imgData, err := io.ReadAll(fileReader)

			imgPathParts := strings.Split(f.Name, "/")
			filename := imgPathParts[len(imgPathParts)-1] // includes extension

			newImgPath := fmt.Sprintf("%s/%s", a.tmpDirectoryName, filename)

			err = os.WriteFile(newImgPath, imgData, 0666)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// for now, just show the first html file
	runtime.EventsEmit(a.ctx, "page", string(a.htmlFiles[0]))
}

func (a *App) UpdatePage() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "page", string(a.htmlFiles[a.currPage]))
	}
}

func (a *App) PrevFile() {
	if len(a.htmlFiles) == 0 {
		return
	}

	a.currPage = a.currPage - 1
	if a.currPage < 0 {
		a.currPage = 0
	}

	a.UpdatePage()
}

func (a *App) NextFile() {
	if len(a.htmlFiles) == 0 {
		return
	}

	a.currPage = a.currPage + 1
	if a.currPage > len(a.htmlFiles)-1 {
		a.currPage = 0
	}

	a.UpdatePage()
}
