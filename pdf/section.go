package pdf

import "github.com/AlexNa-Holdings/savva-reports/cmn"

func (doc *Doc) NewSection(title string) {
	// Add a new page for the section
	doc.NextPage()

	doc.SetFont("DejaVuBold", "", 24)
	doc.SetTextColor(0, 0, 0) // Black

	doc.TextCentered(title, cmn.PageWidth/2, doc.cy)

	// Draw a line under the title
	doc.SetLineWidth(1)
	doc.Line(doc.margin_left, 70, doc.pageWidth-doc.margin_right, 70)

	doc.cy = 70

	doc.NewNLines(10)

}
