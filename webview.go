
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io"
	"net/http"

	"time"
)

type WebView struct {
	content        *widget.RichText
	scroll         *container.Scroll
	onURLChange    func(string)
	onLoadComplete func(string, error)
	currentURL     string
}

func NewWebView(onURLChange func(string), onLoadComplete func(string, error)) *WebView {
	wv := &WebView{
		onURLChange:    onURLChange,
		onLoadComplete: onLoadComplete,
	}

	wv.content = widget.NewRichTextFromMarkdown("")
	wv.content.Wrapping = fyne.TextWrapWord
	wv.scroll = container.NewScroll(wv.content)

	return wv
}

func (wv *WebView) Container() fyne.CanvasObject {
	return wv.scroll
}

func (wv *WebView) LoadURL(url string) {
	wv.currentURL = url
	if wv.onURLChange != nil {
		wv.onURLChange(url)
	}

	wv.content.ParseMarkdown("_Loading..._")

	go func() {
		html, finalURL, err := fetchHTML(url)
		if err != nil {
			wv.content.ParseMarkdown(fmt.Sprintf("# Error Loading Page\n\n**%v**\n\nPlease check the URL and try again.", err))
			if wv.onLoadComplete != nil {
				wv.onLoadComplete(url, err)
			}
			return
		}

		// Render HTML as formatted text
		rendered := renderHTML(html, finalURL)
		wv.content.ParseMarkdown(rendered)
		wv.currentURL = finalURL

		if wv.onURLChange != nil {
			wv.onURLChange(finalURL)
		}
		if wv.onLoadComplete != nil {
			wv.onLoadComplete(finalURL, nil)
		}
	}()
}

func (wv *WebView) Reload() {
	if wv.currentURL != "" {
		wv.LoadURL(wv.currentURL)
	}
}

func fetchHTML(url string) (string, string, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", "", fmt.Errorf("failed to fetch: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), resp.Request.URL.String(), nil
}
