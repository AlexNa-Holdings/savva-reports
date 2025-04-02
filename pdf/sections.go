package pdf

import (
	"github.com/AlexNa-Holdings/savva-reports/assets"
	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

func (doc *Doc) NewSection(title string) {
	// Add a new page for the section
	doc.NextPage()

	// make sure it is even
	if doc.CurentPage&1 == 0 {
		doc.NextPage()
	}

	doc.Sections = append(doc.Sections, Section{Title: title, Page: doc.CurentPage})
	doc.Section = len(doc.Sections) - 1
	doc.SubSection = 0
	doc.SubSubSection = 0

	// print cmn.CoverImg on all page
	if err := doc.ImageFrom(assets.PageBgImg, 0, 0, &gopdf.Rect{W: cmn.PageWidth, H: cmn.PageHeight}); err != nil {
		log.Error().Err(err).Msg("Failed to draw cover image")
	}

	doc.SetFont("DejaVuBold", "", 24)
	doc.SetTextColor(0, 0, 0) // Black

	doc.TextCentered(title, 0, doc.GetY())

	doc.SetY(doc.GetY() + 20) // Add some space below the title

	// Draw a line under the title
	doc.SetLineWidth(1)
	doc.Line(doc.Margins.Left, doc.GetY(), doc.PageWidth-doc.Margins.Right, doc.GetY())

	doc.SetY(doc.GetY() + 60)
	doc.SetX(doc.Margins.Left)

}

func (doc *Doc) NewSubSection(title string) {

	if doc.Section == -1 {
		log.Error().Msg("No section started")
		return // No section started
	}

	doc.NewNLines(2)

	if doc.GetY() > doc.PageHeight-doc.Margins.Bottom-100 {
		doc.NextPage()
	}

	doc.Sections[doc.Section].SubSections = append(doc.Sections[doc.Section].SubSections, Section{Title: title, Page: doc.CurentPage})
	doc.SubSection = len(doc.Sections[doc.Section].SubSections) - 1
	doc.SubSubSection = 0

	doc.SetX(doc.Margins.Left)
	doc.SetFont("DejaVuBold", "", 18)
	doc.SetTextColor(0, 0, 0) // Black

	doc.Text(title)

	// Draw a line under the title
	doc.SetLineWidth(1)
	doc.SetY(doc.GetY() + 4)
	doc.SetStrokeColor(SAVVA_DARK_COLOR.R, SAVVA_DARK_COLOR.G, SAVVA_DARK_COLOR.B)
	doc.Line(doc.Margins.Left, doc.GetY(), doc.PageWidth-doc.Margins.Right, doc.GetY())

	doc.SetY(doc.GetY() + 40)
	doc.SetX(doc.Margins.Left)
}

func (doc *Doc) NewSubSubSection(title string) {
	if doc.Section == -1 {
		log.Error().Msg("No section started")
		return // No section started
	}

	doc.NewNLines(1)

	if doc.GetY() > doc.PageHeight-doc.Margins.Bottom-100 {
		doc.NextPage()
	}

	doc.Sections[doc.Section].SubSections[doc.SubSection].SubSections = append(doc.Sections[doc.Section].SubSections[doc.SubSection].SubSections, Section{Title: title, Page: doc.CurentPage})
	doc.SubSubSection = len(doc.Sections[doc.Section].SubSections[doc.SubSection].SubSections) - 1

	doc.SetX(doc.Margins.Left + 20)
	doc.SetFont("TimesBold", "", 16)
	doc.SetTextColor(0, 0, 0) // Black

	doc.Text(title)

	// Draw a line under the title
	doc.SetLineWidth(1)
	doc.SetY(doc.GetY() + 4)
	// doc.SetStrokeColor(SAVVA_DARK_COLOR.R, SAVVA_DARK_COLOR.G, SAVVA_DARK_COLOR.B)
	// doc.Line(doc.margins.Left+20, doc.GetY(), doc.pageWidth-doc.margins.Right, doc.GetY())

	doc.SetY(doc.GetY() + 30)
	doc.SetX(doc.Margins.Left)
}
