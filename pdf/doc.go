package pdf

import (
	"image"

	"github.com/AlexNa-Holdings/savva-reports/assets"
	"github.com/AlexNa-Holdings/savva-reports/data"
	"github.com/AlexNa-Holdings/savva-reports/i18n"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

type PaddingDescription struct {
	Left, Top, Right, Bottom float64
}

type Style struct {
	FontName  string
	FontSize  float64
	FontColor *Color
	BGColor   *Color
	Padding   PaddingDescription
	Align     uint8 // 'L', 'C', 'R'
	MinHeight float64
}

type Color struct {
	R, G, B uint8
}

var SAVVA_COLOR Color = Color{0xff, 0x71, 0x00}
var SAVVA_DARK_COLOR Color = Color{0xc4, 0x80, 0x00}

type Section struct {
	Title       string
	Page        int
	SubSections []*Section
}

type Doc struct { // Extended gopdf.GoPdf
	*gopdf.GoPdf
	Locale      string
	CurentPage  int
	UserAddress string
	Sections    []*Section

	Margins                            PaddingDescription
	PageWidth, PageHeight              float64
	Section, SubSection, SubSubSection int
	PrintHeader                        bool

	// data
	History   []data.HistoryRecord
	Sponsored []data.Sponsored

	// internal
	indent       int
	indentWidth  float64
	skip_newline bool
	style        Style
	styles       []Style
	GetImage     func(string) (image.Image, error)
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
		Section:     -1,
		SubSection:  -1,
		PrintHeader: true,
	}

	doc.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	doc.SetMargins(0, 0, 0, 0) // No margins

	doc.PageWidth, doc.PageHeight = gopdf.PageSizeA4.W, gopdf.PageSizeA4.H

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

func (doc *Doc) SetDocFont(fontName string, size float64) {
	if err := doc.SetFont(fontName, "", size); err != nil {
		log.Error().Err(err).Msgf("Failed to set font %s", fontName)
	}
	doc.style.FontName = fontName
	doc.style.FontSize = size
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

	doc.SetDocFont(doc.style.FontName, doc.style.FontSize)
	doc.SetTextColor(
		doc.style.FontColor.R,
		doc.style.FontColor.G,
		doc.style.FontColor.B,
	)
}

func (doc *Doc) GetMarginWidth() float64 {
	return doc.PageWidth - doc.Margins.Left - doc.Margins.Right
}

func (doc *Doc) GetMarginHeight() float64 {
	return doc.PageHeight - doc.Margins.Top - doc.Margins.Bottom
}

func (doc *Doc) AssureVertialSpace(h float64) {
	if doc.GetY()+h > doc.PageHeight-doc.Margins.Bottom {
		doc.NextPage()
	}
}
