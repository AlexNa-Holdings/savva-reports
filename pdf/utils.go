package pdf

import "github.com/AlexNa-Holdings/savva-reports/cmn"

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
