package reports

import (
	"math/big"
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

	t.ColStyle[1].Align = 'R'
	t.ColStyle[2].Align = 'R'

	var v *big.Int = new(big.Int)
	var fiat float64

	v = calcSumIn(doc, "token", "transfer")

	log.Debug().Msgf("calcSumIn: %s", v.String())

	fiat = pdf.Value2Float(v, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.savva_in"), doc.FormatValue(v, 18), doc.FormatFiat(fiat))
	doc.WriteTable(t)

}

func calcSumIn(doc *pdf.Doc, contract, event string) *big.Int {
	sum := new(big.Int)

	if doc.History == nil {
		log.Error().Msg("History is nil")
		return sum
	}

	for _, h := range doc.History {
		if contract != "" && h.Contract.String != contract {
			continue
		}

		if event != "" && h.Type.String != event {
			continue
		}

		if doc.UserAddress != h.ToAddr.String {
			continue
		}

		if h.Amount != nil {
			sum = sum.Add(sum, h.Amount)
		}
	}

	return sum
}
