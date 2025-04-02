package pdf

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/signintech/gopdf"
)

func (doc *Doc) TextCentered(text string, x, y float64) {
	// Center the text horizontally

	if x == 0 { // use center of margins
		x = (doc.PageWidth-doc.Margins.Left-doc.Margins.Right)/2 + doc.Margins.Left
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
	if style.FontName != "" && style.FontSize != 0 {
		doc.SetDocFont(style.FontName, style.FontSize)
	}

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

	doc.Margins.Top = 80
	doc.Margins.Bottom = 60

	if doc.CurentPage&1 == 1 {
		doc.Margins.Left = 40
		doc.Margins.Right = 60
	} else {
		doc.Margins.Left = 60
		doc.Margins.Right = 40
	}

	if doc.PrintHeader {
		doc.Header()
	}
	doc.Footer()

	doc.SetX(doc.Margins.Left)
	doc.SetY(doc.Margins.Top)
	doc.NewLine()

}

func (doc *Doc) Header() {
	// Draw page number
	text := fmt.Sprintf("%d", doc.CurentPage)
	doc.SetTextColor(0, 0, 0) // Black
	doc.SetFont("Arial", "", 12)

	if doc.CurentPage&1 == 1 {
		doc.TextRight(text, doc.PageWidth-doc.Margins.Right+20, doc.Margins.Top-20)
	} else {
		doc.TextLeft(text, doc.Margins.Left-20, doc.Margins.Top-20)
	}
}
func (doc *Doc) Footer() {
	// print date of generation
	doc.SetFont("Mono", "", 10)
	doc.TextCentered(fmt.Sprintf("Generated on: %s",
		time.Now().UTC().Format(time.RFC822)),
		0,
		cmn.PageHeight-doc.Margins.Bottom+10)
}

func (doc *Doc) NewLine() {
	doc.SetY(doc.GetY() + doc.style.FontSize*1.3)
	doc.SetX(doc.Margins.Left + float64(doc.indent)*doc.indentWidth)
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
	doc.SetDocFont(style.FontName, style.FontSize)

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
	doc.SetDocFont(style.FontName, style.FontSize)
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
			if doc.GetY() > doc.PageHeight-doc.Margins.Bottom {
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
		if i > 0 && (n-i)%3 == 0 && intPart[i-1] != '-' {
			intPartWithCommas += ","
		}
		intPartWithCommas += string(ch)
	}

	if fracPart != "" {
		return intPartWithCommas + "." + fracPart
	}
	return intPartWithCommas
}

// addCommas adds commas to a numeric string.
func addCommas(number string) string {
	var result strings.Builder
	n := len(number)

	for i, digit := range number {
		if (n-i)%3 == 0 && i != 0 && number[i-1] != '-' {
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
		str = strings.ReplaceAll(str, ",", " ")
		str = strings.ReplaceAll(str, ".", ",")
	}

	return str
}

func (doc *Doc) FormatFiat(value float64) string {
	str := fmt.Sprintf("%.2f", value)

	p := strings.Split(str, ".")

	str = addCommas(p[0]) + "." + p[1]

	if doc.Locale == "ru" {
		//replace . -> ' ', . -> ,
		str = strings.ReplaceAll(str, ",", " ")
		str = strings.ReplaceAll(str, ".", ",")
	}

	return cmn.C.CurrencySymbol + str
}

// DrawImageCover crops and scales the image to cover the given area.
func (doc *Doc) DrawImageCover(img image.Image, x, y, targetW, targetH float64) error {
	bounds := img.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	// Calculate aspect ratios
	targetRatio := targetW / targetH
	imageRatio := imgW / imgH

	var cropRect image.Rectangle

	if imageRatio > targetRatio {
		// Image is wider than target: crop horizontally
		newWidth := int(targetRatio * imgH)
		startX := (bounds.Dx() - newWidth) / 2
		cropRect = image.Rect(startX, 0, startX+newWidth, bounds.Dy())
	} else {
		// Image is taller than target: crop vertically
		newHeight := int(imgW / targetRatio)
		startY := (bounds.Dy() - newHeight) / 2
		cropRect = image.Rect(0, startY, bounds.Dx(), startY+newHeight)
	}

	// Crop the image
	croppedImg := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))
	draw.Draw(croppedImg, croppedImg.Bounds(), img, cropRect.Min, draw.Src)

	// Draw it scaled into the target rect
	return doc.ImageFrom(croppedImg, x, y, &gopdf.Rect{W: targetW, H: targetH})
}

func (doc *Doc) DrawBigImage(img image.Image) error {
	img = EnsureRGBA(img)

	bounds := img.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	// Maximum drawable width and height within margins
	maxW := doc.PageWidth - doc.Margins.Left - doc.Margins.Right
	maxH := doc.PageHeight - doc.Margins.Bottom - doc.Margins.Top

	// Calculate scale factor (only shrink)
	scaleX := maxW / imgW
	scaleY := maxH / imgH
	scale := math.Min(scaleX, scaleY)

	// Never enlarge
	if scale > 1 {
		scale = 1
	}

	// Final dimensions
	drawW := imgW * scale
	drawH := imgH * scale

	// Check if there's enough space on this page
	if doc.GetY()+drawH > doc.PageHeight-doc.Margins.Bottom {
		doc.NextPage()
	}

	// Center horizontally
	x := doc.Margins.Left + (maxW-drawW)/2
	y := doc.GetY()

	// Draw the image
	err := doc.ImageFrom(img, x, y, &gopdf.Rect{W: drawW, H: drawH})
	if err != nil {
		return err
	}

	// Move cursor below the image
	doc.SetY(doc.GetY() + drawH + 5)
	return nil
}

// EnsureRGBA checks if the image is already an 8-bit RGBA-compatible image,
// and converts it only if needed.
func EnsureRGBA(src image.Image) *image.RGBA {
	// If already *image.RGBA, return as-is
	if rgba, ok := src.(*image.RGBA); ok {
		return rgba
	}

	// If already *image.NRGBA (which Go also handles fine), convert safely
	if nrgba, ok := src.(*image.NRGBA); ok {
		bounds := nrgba.Bounds()
		dst := image.NewRGBA(bounds)
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				dst.Set(x, y, nrgba.At(x, y))
			}
		}
		return dst
	}

	// If other formats (possibly 16-bit or unsupported), convert pixel-by-pixel
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			dst.Set(x, y, color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: uint8(a >> 8),
			})
		}
	}
	return dst
}

func (doc *Doc) EclipseToWidth(s string, width float64) (string, float64) {
	// Check if the string fits within the given width
	textWidth, _ := doc.MeasureTextWidth(s)
	if textWidth <= width {
		return s, textWidth
	}

	// If it doesn't fit, truncate and add ellipsis
	ellipsis := "..."
	truncated := s

	for len(truncated) > 0 {
		truncated = truncated[:len(truncated)-1]
		textWidth, _ = doc.MeasureTextWidth(truncated + ellipsis)
		if textWidth <= width {
			break
		}
	}

	return truncated + ellipsis, textWidth
}

func (doc *Doc) EclipseToWidthWithStyle(s string, width float64, style *Style) (string, float64) {
	doc.SetDocFont(style.FontName, style.FontSize)
	return doc.EclipseToWidth(s, width)
}
