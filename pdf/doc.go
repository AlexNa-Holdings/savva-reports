package pdf

import (
	"github.com/AlexNa-Holdings/savva-reports/assets"
	"github.com/AlexNa-Holdings/savva-reports/data"
	"github.com/AlexNa-Holdings/savva-reports/i18n"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

type Padding struct {
	left, top, right, bottom float64
}

type Style struct {
	FontName  string
	FontSize  float64
	FontColor *Color
	BGColor   *Color
	Padding   Padding
	Align     uint8 // 'L', 'C', 'R'
}

type Color struct {
	R, G, B uint8
}

var SAVVA_COLOR Color = Color{0xff, 0x71, 0x00}
var SAVVA_DARK_COLOR Color = Color{0xc4, 0x80, 0x00}

type Doc struct { // Extended gopdf.GoPdf
	*gopdf.GoPdf
	Locale                string
	CurentPage            int
	UserAddress           string
	margin_left           float64
	margin_right          float64
	margin_top            float64
	margin_bottom         float64
	pageWidth, pageHeight float64
	indent                int
	indentWidth           float64
	skip_newline          bool
	style                 Style
	styles                []Style

	// data
	History []data.HistoryRecord
}

func NewDoc(user_addr, locale string) (*Doc, error) {
	// Create a new PDF document.
	doc := Doc{
		GoPdf:       new(gopdf.GoPdf),
		UserAddress: user_addr,
		Locale:      locale,
		CurentPage:  0,
		style: Style{
			FontName:  "Arial",
			FontSize:  12,
			FontColor: &Color{0, 0, 0},
		},
		indentWidth: 20,
	}
	doc.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	doc.pageWidth, doc.pageHeight = gopdf.PageSizeA4.W, gopdf.PageSizeA4.H

	// Load all fonts
	for name, font := range assets.AllFonts {
		err := doc.AddTTFFontData(name, font)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to load font %s", name)
			return nil, err
		}
	}

	// Set default font.  Handle the error.
	if err := doc.SetFont("Arial", "", 24); err != nil {
		log.Error().Err(err).Msg("Failed to set font for cover page")
		return nil, err
	}
	doc.SetTextColor(0, 0, 0) // Black

	return &doc, nil
}

func (doc *Doc) T(key string) string {
	return i18n.T(key, doc.Locale)
}

func (doc *Doc) setFont(fontName string, size float64) {
	if err := doc.SetFont(fontName, "", size); err != nil {
		log.Error().Err(err).Msgf("Failed to set font %s", fontName)
	}
	doc.style.FontName = fontName
	doc.style.FontSize = float64(size)
}

func (doc *Doc) SetColor(c *Color) {
	doc.SetTextColor(c.R, c.G, c.B)
	doc.style.FontColor = c
}

func (doc *Doc) saveStyle() {
	doc.styles = append(doc.styles, doc.style)
}

func (doc *Doc) restoreStyle() {
	if len(doc.styles) == 0 {
		log.Error().Msg("No styles to restore")
		return
	}
	doc.style = doc.styles[len(doc.styles)-1]
	doc.styles = doc.styles[:len(doc.styles)-1]

	doc.setFont(doc.style.FontName, doc.style.FontSize)
	doc.SetTextColor(
		doc.style.FontColor.R,
		doc.style.FontColor.G,
		doc.style.FontColor.B,
	)
}

func (doc *Doc) GetMarginWidth() float64 {
	return doc.pageWidth - doc.margin_left - doc.margin_right
}
func (doc *Doc) GetMarginHeight() float64 {
	return doc.pageHeight - doc.margin_top - doc.margin_bottom
}
