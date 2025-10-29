

// browser.go
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

type Browser struct {
	window      fyne.Window
	urlEntry    *widget.Entry
	webview     *WebView
	statusLabel *widget.Label
	backBtn     *widget.Button
	forwardBtn  *widget.Button
	refreshBtn  *widget.Button
	homeBtn     *widget.Button
	history     *History
}

func NewBrowser(w fyne.Window) *Browser {
	b := &Browser{
		window:  w,
		history: NewHistory(),
	}

	b.webview = NewWebView(b.onURLChange, b.onLoadComplete)
	b.setupUI()

	return b
}

func (b *Browser) setupUI() {
	// Navigation buttons
	b.backBtn = widget.NewButton("←", b.goBack)
	b.backBtn.Importance = widget.LowImportance
	
	b.forwardBtn = widget.NewButton("→", b.goForward)
	b.forwardBtn.Importance = widget.LowImportance
	
	b.refreshBtn = widget.NewButton("⟳", b.refresh)
	b.refreshBtn.Importance = widget.LowImportance
	
	b.homeBtn = widget.NewButton("⌂", b.goHome)
	b.homeBtn.Importance = widget.LowImportance

	// URL entry
	b.urlEntry = widget.NewEntry()
	b.urlEntry.SetPlaceHolder("Enter URL (e.g. https://example.com)")
	b.urlEntry.OnSubmitted = func(url string) {
		b.navigate(url)
	}

	// Status bar
	b.statusLabel = widget.NewLabel("Ready")
	b.statusLabel.TextStyle.Italic = true

	b.updateNavButtons()
}

func (b *Browser) BuildUI() fyne.CanvasObject {
	// Go button
	goBtn := widget.NewButton("Go", func() {
		b.navigate(b.urlEntry.Text)
	})
	goBtn.Importance = widget.HighImportance

	// Navigation bar
	navButtons := container.NewHBox(
		b.backBtn,
		b.forwardBtn,
		b.refreshBtn,
		b.homeBtn,
	)

	urlBar := container.NewBorder(nil, nil, navButtons, goBtn, b.urlEntry)

	// Main layout
	return container.NewBorder(
		urlBar,
		b.statusLabel,
		nil,
		nil,
		b.webview.Container(),
	)
}

func (b *Browser) navigate(urlStr string) {
	if urlStr == "" {
		return
	}

	// Add https:// if no protocol specified
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	b.statusLabel.SetText("Loading...")
	b.history.Add(urlStr)
	b.webview.LoadURL(urlStr)
	b.updateNavButtons()
}

func (b *Browser) goBack() {
	if url := b.history.Back(); url != "" {
		b.webview.LoadURL(url)
		b.updateNavButtons()
	}
}

func (b *Browser) goForward() {
	if url := b.history.Forward(); url != "" {
		b.webview.LoadURL(url)
		b.updateNavButtons()
	}
}

func (b *Browser) refresh() {
	b.webview.Reload()
}

func (b *Browser) goHome() {
	b.navigate("https://www.google.com")
}

func (b *Browser) updateNavButtons() {
	if b.history.CanGoBack() {
		b.backBtn.Enable()
	} else {
		b.backBtn.Disable()
	}
	
	if b.history.CanGoForward() {
		b.forwardBtn.Enable()
	} else {
		b.forwardBtn.Disable()
	}
}

func (b *Browser) onURLChange(url string) {
	b.urlEntry.SetText(url)
}

func (b *Browser) onLoadComplete(url string, err error) {
	if err != nil {
		b.statusLabel.SetText("Error: " + err.Error())
	} else {
		b.statusLabel.SetText("Loaded: " + url)
	}
}

