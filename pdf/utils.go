package pdf

import (
	"fmt"
	"math"
	"math/big"
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

func (doc *Doc) TextWidthStyle(text string, x, y, w float64, style *Style) {
	doc.setFont(style.FontName, style.FontSize)

	switch style.Align {
	case 'L':
		doc.TextLeft(text, x+style.Padding.Left, y)
	case 'C':
		doc.TextCentered(text, x+style.Padding.Left+w/2, y)
	case 'R':
		doc.TextRight(text, x+w-style.Padding.Right, y)
	default:
		doc.TextLeft(text, x+style.Padding.Left, y)
	}
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

	doc.Header()
	doc.Footer()

	doc.SetX(doc.margin_left)
	doc.SetY(doc.margin_top)
	doc.NewLine()

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
	doc.SetY(doc.GetY() + doc.style.FontSize*1.3)
	doc.SetX(doc.margin_left + float64(doc.indent)*doc.indentWidth)

	// Check page boundary
	if doc.GetY() > doc.pageHeight-doc.margin_bottom {
		doc.NextPage()
	}

	log.Debug().Msgf("NewLine: cx=%f, cy=%f", doc.GetX(), doc.GetY())
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

	width = width - style.Padding.Left - style.Padding.Right

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

	return lineHeight*float64(lines) + style.Padding.Top + style.Padding.Bottom
}

// writeTextInWidth writes text at the x, y within the specified width, using provided style.
func (doc *Doc) writeTextInWidth(text string, x, y, width float64, style *Style) {

	if style == nil {
		style = &doc.style
	}

	// Save original style to restore later.
	doc.saveStyle()
	defer doc.restoreStyle()

	width -= style.Padding.Left + style.Padding.Right

	// Set desired style for text
	doc.setFont(style.FontName, style.FontSize)
	doc.SetTextColor(style.FontColor.R, style.FontColor.G, style.FontColor.B)

	text_height, _ := doc.MeasureCellHeightByText("A")
	y += text_height + style.Padding.Top

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
			switch style.Align {
			case 'L':
				doc.TextLeft(line, x+style.Padding.Left, y)
			case 'C':
				doc.TextCentered(line, x+style.Padding.Left+width/2, y)
			case 'R':
				doc.TextRight(line, x+width-style.Padding.Right, y)
			default:
				doc.TextLeft(line, x+style.Padding.Left, y)
			}

			line = word
			lineWidth = wordWidth

			// Move to next line
			y += lineHeight

			// Page break check
			if doc.GetY() > doc.pageHeight-doc.margin_bottom {
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
		switch style.Align {
		case 'L':
			doc.TextLeft(line, x+style.Padding.Left, y)
		case 'C':
			doc.TextCentered(line, x+style.Padding.Left+width/2, y)
		case 'R':
			doc.TextRight(line, x+width-style.Padding.Right, y)
		default:
			doc.TextLeft(line, x+style.Padding.Left, y)
		}
	}
}

func Value2Float(v *big.Int, decimals int) float64 {
	if v == nil || v.Cmp(big.NewInt(0)) == 0 {
		return 0
	}

	// Convert v to big.Float
	f := new(big.Float).SetInt(v)

	// Calculate divisor = 10^decimals as big.Float
	divisor := new(big.Float).SetFloat64(math.Pow(10, float64(decimals)))

	// Divide value by divisor
	f.Quo(f, divisor)

	// Convert to float64
	value, _ := f.Float64()

	return value
}

func FormatValue(amount *big.Int, decimals int) string {
	var precision = 2
	str := amount.String()

	if len(str) <= decimals {
		str = strings.Repeat("0", decimals-len(str)+1) + str
	}

	decimal := str[len(str)-decimals:]
	// make sure we exactly 'precision" digits after the decimal point
	if len(decimal) > precision {
		decimal = decimal[:precision]
	} else if len(decimal) < precision {
		decimal = decimal + strings.Repeat("0", precision-len(decimal))
	}
	str = str[:len(str)-decimals] + "." + decimal

	// // Trim trailing zeros and the decimal point if necessary
	// str = strings.TrimRight(str, "0")
	// str = strings.TrimRight(str, ".")

	// Add commas to the integer part
	parts := strings.Split(str, ".")
	intPart := parts[0]
	fracPart := ""
	if len(parts) > 1 {
		fracPart = parts[1]
	}

	intPartWithCommas := ""
	n := len(intPart)
	for i, ch := range intPart {
		if i > 0 && (n-i)%3 == 0 {
			intPartWithCommas += ","
		}
		intPartWithCommas += string(ch)
	}

	if fracPart != "" {
		return intPartWithCommas + "." + fracPart
	}
	return intPartWithCommas
}

func FormatFloat(value float64) string {
	switch {
	case value == 0:
		return "0"
	case value < 0.01:
		return fmt.Sprintf("%.4f", value)
	case value < 1:
		return fmt.Sprintf("%.2f", value)
	default:
		return addCommas(fmt.Sprintf("%.0f", value))
	}
}

// addCommas adds commas to a numeric string.
func addCommas(number string) string {
	var result strings.Builder
	n := len(number)

	for i, digit := range number {
		if (n-i)%3 == 0 && i != 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}

	return result.String()
}

func (doc *Doc) FormatValue(amount *big.Int, decimals int) string {
	str := FormatValue(amount, decimals)

	if doc.Locale == "ru" {
		//replace . -> ' ', . -> ,
		str = strings.ReplaceAll(str, ".", " ")
		str = strings.ReplaceAll(str, ",", ".")
	}

	return str
}

func (doc *Doc) FormatFloat(value float64) string {
	str := FormatFloat(value)

	if doc.Locale == "ru" {
		//replace . -> ' ', . -> ,
		str = strings.ReplaceAll(str, ".", " ")
		str = strings.ReplaceAll(str, ",", ".")
	}
	return str
}

func (doc *Doc) FormatFiat(value float64) string {
	str := fmt.Sprintf("%.2f", value)

	if doc.Locale == "ru" {
		//replace . -> ' ', . -> ,
		str = strings.ReplaceAll(str, ".", " ")
		str = strings.ReplaceAll(str, ",", ".")
	}

	return cmn.C.CurrencySymbol + str
}
