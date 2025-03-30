package pdf

import (
	"bytes"

	"github.com/rs/zerolog/log"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
)

// MarkDownToPdf converts Markdown text into PDF content.
func (doc *Doc) MarkDownToPdf(md string) error {
	htmlData := blackfriday.Run([]byte(md))
	node, err := html.Parse(bytes.NewReader(htmlData))
	if err != nil {
		return err
	}
	doc.renderNode(node)

	return nil
}

func (doc *Doc) renderNode(n *html.Node) {
	if n == nil {
		return
	}

	switch n.Type {
	case html.ElementNode:
		doc.handleElementStart(n)
	case html.TextNode:

		log.Debug().Msgf("Text: (%s)", n.Data)

		// text := strings.TrimSpace(n.Data)
		text := n.Data
		if text != "" {
			doc.writeText(text)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		doc.renderNode(child)
	}

	if n.Type == html.ElementNode {
		doc.handleElementEnd(n)
	}
}

func (doc *Doc) handleElementStart(n *html.Node) {

	log.Debug().Msgf("< %s", n.Data)

	doc.saveStyle()

	switch n.Data {
	case "h1":
		doc.setFont("Times", 24)
		doc.NewLine()
	case "h2":
		doc.setFont("Times", 20)
		doc.NewLine()
	case "h3":
		doc.setFont("Times", 18)
		doc.NewLine()
	case "p":
		doc.setFont("Times", 14)
		if !doc.skip_newline {
			doc.NewLine()
			doc.skip_newline = false
		}

	case "strong":
		doc.setFont("TimesBold", 14)
		doc.SetColor(SAVVA_DARK_COLOR)
	case "em":
		doc.setFont("TimesBold", 14)
	case "ul", "ol":
		doc.indent++
	case "li":
		doc.NewLine()
		doc.writeText("• ")
		doc.skip_newline = true
	case "br":
		doc.NewLine()
	}
}

func (doc *Doc) handleElementEnd(n *html.Node) {
	log.Debug().Msgf("> %s", n.Data)

	switch n.Data {
	case "strong", "em":
		doc.setFont("Times", 12)
	case "h1", "h2", "h3", "p":
		doc.NewLine()
	case "ul", "ol":
		doc.indent--
	}

	doc.restoreStyle()
}
