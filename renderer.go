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
	ctx := &renderContext{
		result:     &result,
		inPre:      false,
		listDepth:  0,
		orderedNum: make(map[int]int),
	}
	
	renderNode(doc, ctx)

	// Clean up the output
	output := result.String()
	output = strings.TrimSpace(output)
	// Remove excessive blank lines (more than 2 consecutive)
	for strings.Contains(output, "\n\n\n\n") {
		output = strings.ReplaceAll(output, "\n\n\n\n", "\n\n\n")
	}
	
	return output
}

type renderContext struct {
	result     *strings.Builder
	inPre      bool
	listDepth  int
	orderedNum map[int]int
	skipText   bool
}

func renderNode(n *html.Node, ctx *renderContext) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "script", "style", "noscript", "iframe", "svg", "canvas":
			// Skip these elements entirely
			return
			
		case "head":
			// Skip head but process title
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "title" {
					ctx.result.WriteString("# ")
					renderChildren(c, ctx)
					ctx.result.WriteString("\n\n---\n\n")
				}
			}
			return

		case "nav", "header", "footer", "aside":
			// Add subtle separation for structural elements
			ctx.result.WriteString("\n")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n")
			return

		case "h1":
			ctx.result.WriteString("\n\n# ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return
		case "h2":
			ctx.result.WriteString("\n\n## ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return
		case "h3":
			ctx.result.WriteString("\n\n### ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return
		case "h4":
			ctx.result.WriteString("\n\n#### ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return
		case "h5":
			ctx.result.WriteString("\n\n##### ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return
		case "h6":
			ctx.result.WriteString("\n\n###### ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return

		case "p":
			ctx.result.WriteString("\n\n")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n")
			return

		case "div", "section", "article", "main":
			// Add spacing for block elements
			ctx.result.WriteString("\n")
			renderChildren(n, ctx)
			return

		case "br":
			ctx.result.WriteString("  \n") // Two spaces + newline for markdown line break
			return

		case "hr":
			ctx.result.WriteString("\n\n---\n\n")
			return

		case "strong", "b":
			ctx.result.WriteString("**")
			renderChildren(n, ctx)
			ctx.result.WriteString("**")
			return

		case "em", "i":
			ctx.result.WriteString("*")
			renderChildren(n, ctx)
			ctx.result.WriteString("*")
			return

		case "u":
			ctx.result.WriteString("_")
			renderChildren(n, ctx)
			ctx.result.WriteString("_")
			return

		case "code":
			if !ctx.inPre {
				ctx.result.WriteString("`")
				renderChildren(n, ctx)
				ctx.result.WriteString("`")
				return
			}

		case "a":
			href := getAttr(n, "href")
			text := extractText(n)
			
			if text != "" {
				if href != "" && !strings.HasPrefix(href, "#") && !strings.HasPrefix(href, "javascript:") {
					ctx.result.WriteString("[")
					ctx.result.WriteString(text)
					ctx.result.WriteString("](")
					ctx.result.WriteString(href)
					ctx.result.WriteString(")")
				} else {
					ctx.result.WriteString(text)
				}
			}
			return

		case "img":
			alt := getAttr(n, "alt")
			src := getAttr(n, "src")
			if alt == "" {
				alt = "image"
			}
			if src != "" {
				ctx.result.WriteString(fmt.Sprintf("![%s](%s)", alt, src))
			}
			return

		case "ul":
			ctx.result.WriteString("\n")
			ctx.listDepth++
			renderChildren(n, ctx)
			ctx.listDepth--
			ctx.result.WriteString("\n")
			return

		case "ol":
			ctx.result.WriteString("\n")
			ctx.listDepth++
			ctx.orderedNum[ctx.listDepth] = 1
			renderChildren(n, ctx)
			delete(ctx.orderedNum, ctx.listDepth)
			ctx.listDepth--
			ctx.result.WriteString("\n")
			return

		case "li":
			ctx.result.WriteString("\n")
			// Add indentation
			for i := 1; i < ctx.listDepth; i++ {
				ctx.result.WriteString("  ")
			}
			// Check if this is an ordered or unordered list
			if num, ok := ctx.orderedNum[ctx.listDepth]; ok {
				ctx.result.WriteString(fmt.Sprintf("%d. ", num))
				ctx.orderedNum[ctx.listDepth]++
			} else {
				ctx.result.WriteString("- ")
			}
			renderChildren(n, ctx)
			return

		case "pre":
			ctx.result.WriteString("\n\n```")
			oldInPre := ctx.inPre
			ctx.inPre = true
			ctx.result.WriteString("\n")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n```\n\n")
			ctx.inPre = oldInPre
			return

		case "blockquote":
			ctx.result.WriteString("\n\n> ")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n\n")
			return

		case "table":
			ctx.result.WriteString("\n\n")
			renderTable(n, ctx)
			ctx.result.WriteString("\n\n")
			return

		case "dl":
			ctx.result.WriteString("\n")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n")
			return

		case "dt":
			ctx.result.WriteString("\n**")
			renderChildren(n, ctx)
			ctx.result.WriteString("**")
			return

		case "dd":
			ctx.result.WriteString("\n  ")
			renderChildren(n, ctx)
			return

		case "span", "label", "time":
			// Inline elements - just render children
			renderChildren(n, ctx)
			return

		case "button", "input", "select", "textarea":
			// Render form elements as text
			text := extractText(n)
			if text != "" {
				ctx.result.WriteString("[")
				ctx.result.WriteString(text)
				ctx.result.WriteString("]")
			}
			value := getAttr(n, "value")
			if value != "" && text != value {
				ctx.result.WriteString(" (")
				ctx.result.WriteString(value)
				ctx.result.WriteString(")")
			}
			return

		case "form":
			ctx.result.WriteString("\n**Form:**\n")
			renderChildren(n, ctx)
			ctx.result.WriteString("\n")
			return
		}
	} else if n.Type == html.TextNode {
		if ctx.skipText {
			return
		}
		
		text := n.Data
		if !ctx.inPre {
			// Collapse whitespace
			text = strings.TrimSpace(text)
			if text != "" {
				// Replace multiple spaces with single space
				text = strings.Join(strings.Fields(text), " ")
				ctx.result.WriteString(text)
				// Add space after text node if not followed by punctuation
				if !endsWithPunctuation(text) {
					ctx.result.WriteString(" ")
				}
			}
		} else {
			ctx.result.WriteString(text)
		}
		return
	}

	// Process children for all other elements
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(c, ctx)
	}
}

func renderChildren(n *html.Node, ctx *renderContext) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(c, ctx)
	}
}

func renderTable(n *html.Node, ctx *renderContext) {
	// Collect table data
	var rows [][]string
	
	for tbody := n.FirstChild; tbody != nil; tbody = tbody.NextSibling {
		if tbody.Type == html.ElementNode {
			if tbody.Data == "thead" || tbody.Data == "tbody" || tbody.Data == "tfoot" {
				for tr := tbody.FirstChild; tr != nil; tr = tr.NextSibling {
					if tr.Type == html.ElementNode && tr.Data == "tr" {
						row := extractTableRow(tr)
						if len(row) > 0 {
							rows = append(rows, row)
						}
					}
				}
			} else if tbody.Data == "tr" {
				row := extractTableRow(tbody)
				if len(row) > 0 {
					rows = append(rows, row)
				}
			}
		}
	}

	if len(rows) == 0 {
		return
	}

	// Find max columns
	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}

	// Render table
	for i, row := range rows {
		ctx.result.WriteString("| ")
		for j := 0; j < maxCols; j++ {
			if j < len(row) {
				ctx.result.WriteString(row[j])
			}
			ctx.result.WriteString(" | ")
		}
		ctx.result.WriteString("\n")

		// Add header separator after first row
		if i == 0 {
			ctx.result.WriteString("|")
			for j := 0; j < maxCols; j++ {
				ctx.result.WriteString(" --- |")
			}
			ctx.result.WriteString("\n")
		}
	}
}

func extractTableRow(tr *html.Node) []string {
	var cells []string
	for td := tr.FirstChild; td != nil; td = td.NextSibling {
		if td.Type == html.ElementNode && (td.Data == "td" || td.Data == "th") {
			text := strings.TrimSpace(extractText(td))
			cells = append(cells, text)
		}
	}
	return cells
}

func extractText(n *html.Node) string {
	var result strings.Builder
	extractTextHelper(n, &result)
	return strings.TrimSpace(result.String())
}

func extractTextHelper(n *html.Node, result *strings.Builder) {
	if n.Type == html.TextNode {
		result.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTextHelper(c, result)
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

func endsWithPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}
	last := s[len(s)-1]
	return last == '.' || last == ',' || last == '!' || last == '?' || last == ';' || last == ':'
}
