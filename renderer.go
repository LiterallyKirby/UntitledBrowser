package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func renderHTML(htmlContent string, baseURL string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return fmt.Sprintf("Error parsing HTML: %v", err)
	}

	var result strings.Builder
	renderNode(doc, &result, 0)

	return result.String()
}

func renderNode(n *html.Node, result *strings.Builder, depth int) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "noscript", "iframe":
			// Skip these elements
			return
		case "h1":
			result.WriteString("\n\n# ")
			renderChildren(n, result, depth)
			result.WriteString("\n")
			return
		case "h2":
			result.WriteString("\n\n## ")
			renderChildren(n, result, depth)
			result.WriteString("\n")
			return
		case "h3":
			result.WriteString("\n\n### ")
			renderChildren(n, result, depth)
			result.WriteString("\n")
			return
		case "h4":
			result.WriteString("\n\n#### ")
			renderChildren(n, result, depth)
			result.WriteString("\n")
			return
		case "p":
			result.WriteString("\n\n")
			renderChildren(n, result, depth)
			return
		case "br":
			result.WriteString("\n")
			return
		case "hr":
			result.WriteString("\n\n---\n\n")
			return
		case "strong", "b":
			result.WriteString("**")
			renderChildren(n, result, depth)
			result.WriteString("**")
			return
		case "em", "i":
			result.WriteString("*")
			renderChildren(n, result, depth)
			result.WriteString("*")
			return
		case "code":
			result.WriteString("`")
			renderChildren(n, result, depth)
			result.WriteString("`")
			return
		case "a":
			href := getAttr(n, "href")
			result.WriteString("[")
			renderChildren(n, result, depth)
			if href != "" {
				result.WriteString("](")
				result.WriteString(href)
				result.WriteString(")")
			} else {
				result.WriteString("]")
			}
			return
		case "img":
			alt := getAttr(n, "alt")
			src := getAttr(n, "src")
			if alt != "" {
				result.WriteString(fmt.Sprintf("![%s](%s)", alt, src))
			} else {
				result.WriteString(fmt.Sprintf("![image](%s)", src))
			}
			return
		case "ul", "ol":
			result.WriteString("\n")
			renderChildren(n, result, depth+1)
			result.WriteString("\n")
			return
		case "li":
			result.WriteString("\n")
			for i := 0; i < depth; i++ {
				result.WriteString("  ")
			}
			result.WriteString("* ")
			renderChildren(n, result, depth)
			return
		case "pre":
			result.WriteString("\n```\n")
			renderChildren(n, result, depth)
			result.WriteString("\n```\n")
			return
		case "blockquote":
			result.WriteString("\n> ")
			renderChildren(n, result, depth)
			result.WriteString("\n")
			return
		}
	} else if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			result.WriteString(text)
			result.WriteString(" ")
		}
		return
	}

	// Process children for all other elements
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(c, result, depth)
	}
}

func renderChildren(n *html.Node, result *strings.Builder, depth int) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(c, result, depth)
	}
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}
