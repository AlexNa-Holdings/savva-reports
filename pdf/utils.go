package pdf

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/rs/zerolog/log"
)

func (doc *Doc) TextCentered(text string, x, y float64) {
	// Center the text horizontally

	if x == 0 { // use center of margins
		x = (doc.pageWidth-doc.margin_left-doc.margin_right)/2 + doc.margin_left
	}
	textWidth, _ := doc.MeasureTextWidth(text)
	doc.SetXY(x-(textWidth/2), y)
	doc.Text(text)
}

func (doc *Doc) TextLeft(text string, x, y float64) {
	// Left-align the text
	doc.SetXY(x, y)
	doc.Text(text)
}

func (doc *Doc) TextRight(text string, x, y float64) {
	// Right-align the text
	textWidth, _ := doc.MeasureTextWidth(text)
	doc.SetXY(x-textWidth, y)
	doc.Text(text)
}

func (doc *Doc) AddBlankPage() {
	doc.AddPage()
	doc.SetTextColor(0xff, 0xff, 0xff)
	doc.TextCentered("empty page", cmn.PageWidth/2, 20)
}

func (doc *Doc) NextPage() {
	doc.AddPage()
	doc.CurentPage++

	doc.margin_top = 80
	doc.margin_bottom = 60

	if doc.CurentPage&1 == 1 {
		doc.margin_left = 40
		doc.margin_right = 60
	} else {
		doc.margin_left = 60
		doc.margin_right = 40
	}

	doc.cx = doc.margin_left
	doc.cy = doc.margin_top

	doc.Header()
	doc.Footer()
}

func (doc *Doc) Header() {
	// Draw page number
	text := fmt.Sprintf("%d", doc.CurentPage)
	doc.SetTextColor(0, 0, 0) // Black
	doc.SetFont("Arial", "", 12)

	if doc.CurentPage&1 == 1 {
		doc.TextRight(text, doc.pageWidth-doc.margin_right+20, doc.margin_top-20)
	} else {
		doc.TextLeft(text, doc.margin_left-20, doc.margin_top-20)
	}
}
func (doc *Doc) Footer() {
	// print date of generation
	doc.SetFont("Mono", "", 10)
	doc.TextCentered(fmt.Sprintf("Generated on: %s",
		time.Now().UTC().Format(time.RFC822)),
		0,
		cmn.PageHeight-doc.margin_bottom+10)
}

func (doc *Doc) NewLine() {
	doc.cy += doc.style.FontSize * 1.3
	doc.cx = doc.margin_left + float64(doc.indent)*doc.indentWidth

	// Check page boundary
	if doc.cy > doc.pageHeight-doc.margin_bottom {
		doc.NextPage()
	}

	log.Debug().Msgf("NewLine: cx=%f, cy=%f", doc.cx, doc.cy)
}

func (doc *Doc) NewNLines(n int) {
	for i := 0; i < n; i++ {
		doc.NewLine()
	}
}

// estimateTextHeight estimates the height required for the provided text within a given width and style.
func (doc *Doc) estimateTextHeight(text string, width float64, style *Style) float64 {
	// Save original document style
	doc.saveStyle()
	defer doc.restoreStyle()

	// Apply temporary style for measurement
	doc.setFont(style.FontName, style.FontSize)

	// Calculate line height (standard is ~1.2x font size)
	lineHeight := style.FontSize * 1.2

	words := strings.Fields(text)
	var currentWidth float64
	var lines int = 1

	spaceWidth, _ := doc.MeasureTextWidth(" ")

	for i, word := range words {
		wordWidth, _ := doc.MeasureTextWidth(word)

		// Add space width if it's not the first word
		if i > 0 {
			wordWidth += spaceWidth
		}

		if currentWidth+wordWidth > float64(width) {
			lines++
			currentWidth = wordWidth - spaceWidth // Reset current line width to current word's width
		} else {
			currentWidth += wordWidth
		}
	}

	return lineHeight * float64(lines)
}

// writeTextInWidth writes text at the current (doc.cx, doc.cy) within the specified width, using provided style.
func (doc *Doc) writeTextInWidth(text string, width float64, style *Style) {
	// Save original style to restore later.
	doc.saveStyle()
	defer doc.restoreStyle()

	// Set desired style for text
	doc.setFont(style.FontName, style.FontSize)
	doc.SetTextColor(style.FontColor.R, style.FontColor.G, style.FontColor.B)

	lineHeight := style.FontSize * 1.2
	spaceWidth, _ := doc.MeasureTextWidth(" ")

	words := strings.Fields(text)
	if len(words) == 0 {
		doc.NewLine()
		return
	}

	var line string
	var lineWidth float64

	for _, word := range words {
		wordWidth, _ := doc.MeasureTextWidth(word)

		// Add space width if line isn't empty
		additionalWidth := wordWidth
		if line != "" {
			additionalWidth += spaceWidth
		}

		if lineWidth+additionalWidth > float64(width) {
			// Write current line to PDF
			doc.writeText(line)
			line = word
			lineWidth = wordWidth

			// Move to next line
			doc.cx = doc.margin_left
			doc.cy += lineHeight

			// Page break check
			if doc.cy > doc.pageHeight-doc.margin_bottom {
				doc.NextPage()
			}
		} else {
			if line != "" {
				line += " "
			}
			line += word
			lineWidth += additionalWidth
		}
	}

	// Write the last remaining line
	if line != "" {
		doc.writeText(line)
		doc.cx += lineWidth
	}
}

func (doc *Doc) writeText(text string) {
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
		maxWidth := doc.pageWidth - doc.margin_right - doc.cx
		line, line_width, remainingWords := doc.wrapText(words, maxWidth)
		words = remainingWords

		// Check for page break
		if doc.cy > doc.pageHeight-doc.margin_bottom {
			doc.NextPage()
		}

		doc.SetXY(doc.cx, doc.cy)
		doc.Text(line)
		if len(remainingWords) > 0 {
			doc.NewLine()
		} else {
			doc.cx += line_width
		}
	}

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
