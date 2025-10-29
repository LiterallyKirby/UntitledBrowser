package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Go Web Browser")
	w.Resize(fyne.NewSize(1200, 800))

	browser := NewBrowser(w)
	w.SetContent(browser.BuildUI())
	
	w.ShowAndRun()
}
