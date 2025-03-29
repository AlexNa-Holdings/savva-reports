package pdf

import "github.com/AlexNa-Holdings/savva-reports/cmn"

func (doc *Doc) NewSection(title string) {
	// Add a new page for the section
	doc.AddPage()
	doc.CurentPage++

	doc.SetFont("Arial", "", 24)
	doc.SetTextColor(0, 0, 0) // Black

	// Center the title on the page
	textWidth, _ := doc.MeasureTextWidth(title)
	doc.SetXY((cmn.PageWidth/2)-(textWidth/2), 50)
	doc.Text(title)

	// Draw a line under the title
	doc.SetLineWidth(1)
	doc.Line(50, 70, cmn.PageWidth-50, 70)
}
