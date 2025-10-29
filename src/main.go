
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

// fetch HTML content
func fetch(url string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

// extract visible text (very simple)
func extractText(r io.Reader) string {
	z := html.NewTokenizer(r)
	var text strings.Builder
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return text.String()
		case html.TextToken:
			t := strings.TrimSpace(string(z.Text()))
			if t != "" {
				text.WriteString(t + "\n")
			}
		}
	}
}

func main() {
	a := app.New()
	w := a.NewWindow("Go Browser")
	w.Resize(fyne.NewSize(800, 600))

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL (e.g. example.com)")
	content := widget.NewMultiLineEntry()
	content.Wrapping = fyne.TextWrapWord
	content.SetMinRowsVisible(25)

	loadBtn := widget.NewButton("Go", func() {
		body, err := fetch(urlEntry.Text)
		if err != nil {
			content.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		text := extractText(strings.NewReader(body))
		content.SetText(text)
	})

	top := container.NewBorder(nil, nil, nil, loadBtn, urlEntry)
	w.SetContent(container.NewBorder(top, nil, nil, nil, content))
	w.ShowAndRun()
}
