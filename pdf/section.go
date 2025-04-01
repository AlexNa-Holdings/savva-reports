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
	doc.Line(doc.margin_left, doc.GetY(), doc.pageWidth-doc.margin_right, doc.GetY())

	doc.SetY(doc.GetY() + 60)
	doc.SetX(doc.margin_left)

}
