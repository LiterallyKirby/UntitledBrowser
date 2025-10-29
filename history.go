package main

type History struct {
	urls    []string
	current int
}

func NewHistory() *History {
	return &History{
		urls:    []string{},
		current: -1,
	}
}

func (h *History) Add(url string) {
	// Remove any forward history
	if h.current < len(h.urls)-1 {
		h.urls = h.urls[:h.current+1]
	}

	// Don't add duplicate of current page
	if h.current >= 0 && h.current < len(h.urls) && h.urls[h.current] == url {
		return
	}

	h.urls = append(h.urls, url)
	h.current = len(h.urls) - 1
}

func (h *History) Back() string {
	if h.CanGoBack() {
		h.current--
		return h.urls[h.current]
	}
	return ""
}

func (h *History) Forward() string {
	if h.CanGoForward() {
		h.current++
		return h.urls[h.current]
	}
	return ""
}

func (h *History) CanGoBack() bool {
	return h.current > 0
}

func (h *History) CanGoForward() bool {
	return h.current < len(h.urls)-1
}

func (h *History) Current() string {
	if h.current >= 0 && h.current < len(h.urls) {
		return h.urls[h.current]
	}
	return ""
}
