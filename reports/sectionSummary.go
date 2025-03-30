package reports

import (
	"time"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/AlexNa-Holdings/savva-reports/data"
	"github.com/AlexNa-Holdings/savva-reports/pdf"
	"github.com/rs/zerolog/log"
)

func addSectionSummary(doc *pdf.Doc, from, to time.Time) {
	// Add a new section for the summary
	doc.NewSection(doc.T("section_summary"))

	if doc.History == nil {
		var err error
		doc.History, err = data.GetHistory(doc.UserAddress, &from, &to)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch history")
			return
		}
	}

	t := pdf.NewTable()
	t.SetHeader(doc.T("description"), "SAVVA", cmn.C.CurrencySymbol)
	t.ColWidths = []float64{0, 100, 100}

	t.AddRow("description", "123", "567")

	doc.WriteTable(t)

}
