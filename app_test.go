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

func TestHtmlRegex(t *testing.T) {
	app := NewApp()
	app.ctx = nil

	htmlRegex := app.HtmlRegex()
	match1 := htmlRegex.MatchString("blah.html")
	if !match1 {
		t.Fatalf("blah.html should be valid")
	}

	match2 := htmlRegex.MatchString("blah.htm")
	if !match2 {
		t.Fatalf("blah.htm should be valid")
	}

	match3 := htmlRegex.MatchString("blah.mklmfkl.html")
	if !match3 {
		t.Fatalf("blah.mklmfkl.html should be valid")
	}

	match4 := htmlRegex.MatchString("blah.xhtml")
	if !match4 {
		t.Fatalf("blah.xhtml should be valid")
	}

	noMatch1 := htmlRegex.MatchString("blah.xhtmlz")
	if noMatch1 {
		t.Fatalf("blah.xhtmlz should not be valid")
	}

	noMatch2 := htmlRegex.MatchString("blah.htmlo")
	if noMatch2 {
		t.Fatalf("blah.htmlo should not be valid")
	}

	noMatch3 := htmlRegex.MatchString("blah.html.sdnfksdnl")
	if noMatch3 {
		t.Fatalf("blah.html.sdnfksdnl should not be valid")
	}

	noMatch4 := htmlRegex.MatchString("blah.hhtml")
	if noMatch4 {
		t.Fatalf("blah.hhtml should not be valid")
	}
}

func TestCssRegex(t *testing.T) {
	app := NewApp()
	app.ctx = nil

	cssRegex := app.CssRegex()

	match1 := cssRegex.MatchString("blah.css")
	if !match1 {
		t.Fatalf("blah.css should be valid")
	}

	noMatch1 := cssRegex.MatchString("blah.cs")
	if noMatch1 {
		t.Fatalf("blah.cs should not be valid")
	}

	noMatch2 := cssRegex.MatchString("blah.css.html")
	if noMatch2 {
		t.Fatalf("blah.css.html should not be valid")
	}

	noMatch3 := cssRegex.MatchString("blah.ccss")
	if noMatch3 {
		t.Fatalf("blah.ccss should not be valid")
	}
}
