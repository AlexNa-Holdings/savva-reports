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

	return doc.MarkDownToPdfEx(md, doc.Margins.Left, doc.GetY(), doc.GetMarginWidth(), doc.GetMarginHeight()-doc.GetY(), true)
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

	doc.renderNode(node, x, y, w, h, auto_page, doc.style.FontSize)
	doc.NewLine()

	return nil
}

func (doc *Doc) renderNode(n *html.Node, x, y, w, h float64, auto_page bool, fontSize float64) (float64, float64, float64, float64) {
	if n == nil {
		return x, y, w, h
	}

	switch n.Type {
	case html.ElementNode:
		doc.handleElementStart(n, x, fontSize)
	case html.TextNode:
		if n.Data != "" {

			log.Debug().Msgf("Text: %s", n.Data)

			x, y, w, h = doc.writeText(n.Data, x, y, w, h, auto_page)
		}
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		x, y, w, h = doc.renderNode(child, x, y, w, h, auto_page, fontSize)
	}

	if n.Type == html.ElementNode {
		doc.handleElementEnd(n, x)
	}

	return x, y, w, h
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
				x = doc.Margins.Left
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
			doc.SetX(x)
		}
	}

	return x, y, w, h
}

func (doc *Doc) handleElementStart(n *html.Node, x, fontSize float64) {

	log.Debug().Msgf("< %s", n.Data)

	doc.saveStyle()

	switch n.Data {
	case "h1":
		doc.setFont("TimesBold", fontSize*1.5)
		if doc.GetX() > x {
			doc.NewLine()
			doc.SetX(x)
		}
	case "h2":
		doc.setFont("TimesBold", fontSize*1.2)
		if doc.GetX() > x {
			doc.NewLine()
			doc.SetX(x)
		}
	case "p":
		doc.setFont("Times", fontSize)
		if !doc.skip_newline {
			doc.NewLine()
			doc.SetX(x)
			doc.skip_newline = false
		}

	case "strong":
		doc.setFont("TimesBold", fontSize)
		doc.SetColor(&SAVVA_DARK_COLOR)
	case "em":
		doc.setFont("TimesBold", fontSize)
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
	case "img":

		if strings.HasPrefix(n.Attr[0].Val, "https://www.youtube.com") ||
			strings.HasPrefix(n.Attr[0].Val, "https://youtube.com") {
			doc.Text("<<YouTube video>> ")
		} else {

			if doc.GetImage != nil {
				img, err := doc.GetImage(n.Attr[0].Val)

				if err != nil {
					log.Error().Err(err).Msgf("Failed to load image %s", n.Attr[0].Val)
				} else {
					doc.DrawBigImage(img)
				}
			}
		}
	case "a":
		doc.setFont("Times", 12.0)
	}
}

func (doc *Doc) handleElementEnd(n *html.Node, x float64) {
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

		if i == 0 && wordWidth > maxWidth {
			// Word is too long to fit in a single line
			// Split the word into smaller parts

			wl := len(word)
			for j := 1; j < wl; j++ {
				w, _ := doc.MeasureTextWidth(word[0 : wl-j])
				if w <= maxWidth {
					line += word[0 : wl-j]
					words[i] = word[wl-j:]
					break
				}
			}

			return line, lineWidth, words[i:]
		}

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
