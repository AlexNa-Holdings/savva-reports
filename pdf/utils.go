package pdf

import (
	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/rs/zerolog/log"
)

func (doc *Doc) TextCentered(text string, x, y float64) {
	// Center the text horizontally
	textWidth, _ := doc.MeasureTextWidth(text)
	doc.SetXY(x-(textWidth/2), y)
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
	// print header & footer

	// set margins

	doc.margin_top = 50
	doc.margin_bottom = 50

	if doc.CurentPage&1 == 0 { // even page
		doc.margin_left = 40
		doc.margin_right = 60
	} else { // odd page
		doc.margin_left = 60
		doc.margin_right = 40
	}

	doc.cx = doc.margin_left
	doc.cy = doc.margin_top
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
