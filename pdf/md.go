package pdf

import (
	"bytes"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
)

// MarkDownToPdf converts Markdown text into PDF content.
func (doc *Doc) MarkDownToPdf(md string) error {
	// htmlData := blackfriday.Run([]byte(md))
	// node, err := html.Parse(bytes.NewReader(htmlData))
	// if err != nil {
	// 	return err
	// }
	// doc.renderNode(node)

	// doc.NewLine()

	return doc.MarkDownToPdfEx(md, doc.margin_left, doc.GetY(), doc.GetMarginWidth(), doc.GetMarginHeight()-doc.GetY(), true)
}

func (doc *Doc) MarkDownToPdfEx(md string, x, y, w, h float64, auto_page bool) error {
	htmlData := blackfriday.Run([]byte(md), blackfriday.WithExtensions(
		blackfriday.CommonExtensions|blackfriday.HardLineBreak,
	))
	node, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return err
	}

	text_height, _ := doc.MeasureCellHeightByText("A")
	doc.SetXY(x, y+text_height)

	doc.renderNode(node, x, y, w, h, auto_page)
	doc.NewLine()

	return nil
}

func (doc *Doc) renderNode(n *html.Node, x, y, w, h float64, auto_page bool) {
	if n == nil {
		return
	}

	switch n.Type {
	case html.ElementNode:
		doc.handleElementStart(n, x)
	case html.TextNode:
		if n.Data != "" {
			x, y, w, h = doc.writeText(n.Data, x, y, w, h, auto_page)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		doc.renderNode(child, x, y, w, h, auto_page)
	}

	if n.Type == html.ElementNode {
		doc.handleElementEnd(n)
	}
}

func (doc *Doc) writeText(text string, x, y, w, h float64, auto_page bool) (float64, float64, float64, float64) {
	words := strings.Fields(text)

	// add leading space
	if len(text) > 0 && text[0] == ' ' {
		words = append([]string{" "}, words...)
	}

	// add trailing space
	if text[len(text)-1] == ' ' {
		words = append(words, " ")
	}

	for len(words) > 0 {
		maxWidth := x + w - doc.GetX()

		line, _, remainingWords := doc.wrapText(words, maxWidth)
		words = remainingWords

		// Check for page break
		if doc.GetY() > y+h {
			if auto_page {
				doc.NextPage()
				x = doc.margin_left
				w = doc.GetMarginWidth()
				h = doc.GetMarginHeight() - doc.GetY()
				return x, y, w, h
			} else {
				return x, y, w, h
			}
		}

		doc.Text(line)
		if len(remainingWords) > 0 {
			doc.NewLine()
		}
	}

	return x, y, w, h
}

func (doc *Doc) handleElementStart(n *html.Node, x float64) {

	log.Debug().Msgf("< %s", n.Data)

	doc.saveStyle()

	switch n.Data {
	case "h1":
		doc.setFont("Times", 24)
		doc.NewLine()
		doc.SetX(x)
	case "h2":
		doc.setFont("Times", 20)
		doc.NewLine()
		doc.SetX(x)
	case "h3":
		doc.setFont("Times", 18)
		doc.NewLine()
		doc.SetX(x)
	case "p":
		doc.setFont("Times", 14)
		if !doc.skip_newline {
			doc.NewLine()
			doc.SetX(x)
			doc.skip_newline = false
		}

	case "strong":
		doc.setFont("TimesBold", 14)
		doc.SetColor(&SAVVA_DARK_COLOR)
	case "em":
		doc.setFont("TimesBold", 14)
	case "ul", "ol":
		doc.indent++
	case "li":
		doc.NewLine()
		doc.SetX(x)
		doc.Text("• ")
		doc.skip_newline = true
	case "br":
		doc.NewLine()
		doc.SetX(x)
	}
}

func (doc *Doc) handleElementEnd(n *html.Node) {
	log.Debug().Msgf("> %s", n.Data)

	switch n.Data {
	case "strong", "em":
		doc.setFont("Times", 12)
	case "h1", "h2", "h3", "p":
		doc.NewLine()
		doc.SetX(x)
	case "ul", "ol":
		doc.indent--
	}

	doc.restoreStyle()
}

func (doc *Doc) wrapText(words []string, maxWidth float64) (string, float64, []string) {
	var line string
	var lineWidth float64
	spaceWidth, _ := doc.MeasureTextWidth(" ")

	for i, word := range words {
		wordWidth, _ := doc.MeasureTextWidth(word)

		// Include space before word if not first word in the line
		additionalWidth := wordWidth
		if i > 0 {
			additionalWidth += spaceWidth
		}

		if lineWidth+additionalWidth > maxWidth {
			return line, lineWidth, words[i:]
		}

		if line == "" {
			line = word
		} else {
			line += " " + word
		}

		lineWidth += additionalWidth
	}

	return line, lineWidth, nil
}
