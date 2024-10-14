package main

import (
	"testing"
)

func TestApp(t *testing.T) {
	app := NewApp()
	htmlFilesLen := len(app.htmlFiles)
	currPage := app.currPage

	if currPage != 0 {
		t.Fatalf("currPage is not 0")
	}

	if htmlFilesLen != 0 {
		t.Fatalf("len of html files is not 0")
	}

	// calling prev file should not change anything
	app.PrevFile()
	if app.currPage != 0 {
		t.Fatalf("currPage is not 0 after calling prevfile")
	}

	// calling next file should not change anything
	app.NextFile()
	if app.currPage != 0 {
		t.Fatalf("currPage is not 0 after calling nextfile")
	}
}

func TestAppNextFile(t *testing.T) {
	app := NewApp()
	app.ctx = nil

	// add some pages
	app.htmlFiles = append(app.htmlFiles, []byte("test1"))
	app.htmlFiles = append(app.htmlFiles, []byte("test2"))
	app.htmlFiles = append(app.htmlFiles, []byte("test3"))

	htmlFilesLen := len(app.htmlFiles)
	if htmlFilesLen != 3 {
		t.Fatalf("len of html files is not 3")
	}

	// calling prev file should not change anything
	app.PrevFile()
	if app.currPage != 0 {
		t.Fatalf("currPage is not 0 after calling prevfile")
	}

	// go to next page
	app.NextFile()
	if app.currPage != 1 {
		t.Fatalf("currPage is not 1 after calling nextfile")
	}

	// go to next page
	app.NextFile()
	if app.currPage != 2 {
		t.Fatalf("currPage is not 2 after calling nextfile")
	}

	// go to prev page
	app.PrevFile()
	if app.currPage != 1 {
		t.Fatalf("currPage is not 1 after calling prevfile")
	}

	// go to last page
	app.NextFile()

	// go to next page (should wraparound back to the first page)
	app.NextFile()
	if app.currPage != 0 {
		t.Fatalf("currPage is not 0 after calling nextfile on last page")
	}
}
