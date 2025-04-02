package reports

import (
	"fmt"

	"github.com/AlexNa-Holdings/savva-reports/pdf"
)

const INDENT = 15.

var TEXT_LEFT = 30.
var TEXT_WIDTH = 200.
var NUMBER_RIGHT = 250.
var NUMBER_WIDTH = 20.

func addTableOfContents(doc *pdf.Doc) {

	doc.PrintHeader = false
	doc.NextPage()

	TEXT_WIDTH = doc.GetMarginWidth() * 2 / 3
	TEXT_LEFT = 20
	NUMBER_RIGHT = doc.GetMarginWidth() - 20

	// make sure it is even
	if doc.CurentPage&1 == 0 {
		doc.NextPage()
	}

	doc.SetDocFont("TimesBold", 24)
	doc.TextCentered(doc.T("table_of_contents"), 0, doc.GetY())

	doc.SetY(doc.GetY() + 20) // Add some space below the title

	// // Draw a line under the title
	// doc.SetLineWidth(1)
	// doc.Line(doc.Margins.Left, doc.GetY(), doc.PageWidth-doc.Margins.Right, doc.GetY())

	doc.SetY(doc.GetY() + 60)
	doc.SetX(doc.Margins.Left)

	indent := 0.

	for _, s := range doc.Sections {
		indent += INDENT
		TOCLine(doc, s, indent, &pdf.Style{
			FontName: "TimesBold",
			FontSize: 14,
			Align:    'L',
		})
		for _, ss := range s.SubSections {
			indent += INDENT
			TOCLine(doc, ss, indent, &pdf.Style{
				FontName: "TimesBold",
				FontSize: 14,
				Align:    'L',
			})
			for _, sss := range ss.SubSections {
				indent += INDENT
				TOCLine(doc, sss, indent, &pdf.Style{
					FontName: "Times",
					FontSize: 14,
					Align:    'L',
				})
				indent -= INDENT
			}
			indent -= INDENT
		}
		indent -= INDENT

	}
}

func TOCLine(doc *pdf.Doc, s *pdf.Section, indent float64, style *pdf.Style) {

	doc.AssureVertialSpace(15)

	t, w := doc.EclipseToWidthWithStyle(s.Title, TEXT_WIDTH-indent, style)

	doc.TextWidthStyle(t, doc.Margins.Left+TEXT_LEFT+indent, doc.GetY(), NUMBER_WIDTH, style)
	doc.TextWidthStyle(fmt.Sprintf("%d", s.Page), doc.Margins.Left+NUMBER_RIGHT, doc.GetY(), NUMBER_WIDTH, &pdf.Style{
		FontName: style.FontName,
		FontSize: style.FontSize,
		Align:    'R',
	})

	// draw the grey line to the number
	doc.SetLineWidth(0.5)
	doc.SetStrokeColor(0xc4, 0xc4, 0xc4)
	doc.Line(doc.Margins.Left+TEXT_LEFT+indent+w+2, doc.GetY(), doc.Margins.Left+NUMBER_RIGHT-2, doc.GetY())

	doc.NewLine()

}
