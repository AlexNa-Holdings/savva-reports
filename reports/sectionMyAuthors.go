package reports

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/AlexNa-Holdings/savva-reports/data"
	"github.com/AlexNa-Holdings/savva-reports/pdf"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

const AVATAR_SIZE = 100.

func addSectionMyAuthors(doc *pdf.Doc, from, to time.Time) {

	if doc.Sponsored == nil {
		var err error
		doc.Sponsored, err = data.GetSponsoredBy(doc.UserAddress)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch sponsored by")
			return
		}
	}

	if len(doc.Sponsored) == 0 {
		return // skip the section
	}

	total := big.NewInt(0)
	for _, s := range doc.Sponsored {
		total.Add(total, s.TotalAmount)
	}

	// Add a new section for the summary
	doc.NewSection(doc.T("section_my_authors"))

	if doc.History == nil {
		var err error
		doc.History, err = data.GetHistory(doc.UserAddress, &from, &to)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch history")
			return
		}
	}

	doc.MarkDownToPdf(fmt.Sprintf(doc.T("my_authors_introduction"), doc.FormatValue(total, 18)) + " SAVVA (" + doc.FormatFiat(pdf.Value2Float(total, 18)*cmn.C.SavvaTokenPrice) + ")")
	doc.NewLine()

	t := pdf.NewTable()
	t.SetHeader(doc.T("account"), "SAVVA", cmn.C.CurrencySymbol)
	t.ColWidths = []float64{0, 100, 100}

	t.ColStyle[0].MinHeight = AVATAR_SIZE
	t.ColStyle[0].Padding.Left = AVATAR_SIZE + 10
	t.ColStyle[1].Align = 'R'
	t.ColStyle[2].Align = 'R'

	t.OnBeforeDrawCell = func(t *pdf.Table, row, col int, x, y float64, w float64, h float64, text string, style *pdf.Style) {
		switch col {
		case 0:
			user, err := data.GetUser(doc.Sponsored[row].Author)
			if err == nil {
				doc.ImageFrom(user.AvatarImg, x+5, y+style.Padding.Top, &gopdf.Rect{W: AVATAR_SIZE, H: AVATAR_SIZE})
			} else {
				log.Error().Err(err).Msg("Failed to fetch user data")
			}
		}
	}

	var fiat float64

	for _, s := range doc.Sponsored {
		fiat = pdf.Value2Float(s.TotalAmount, 18) * cmn.C.SavvaTokenPrice
		user, err := data.GetUser(s.Author)
		if err == nil {

			info := "!MD"
			if user.Name != "" {
				info += "*" + strings.ToUpper(user.Name) + "\u00AE*\n"
			}

			dn := user.GetDisplayName()
			if dn != "" {
				info += dn + "\n"
			}

			info += doc.T("total") + ": \n" + doc.FormatValue(s.TotalAmount, 18) + " SAVVA\n" + doc.FormatFiat(fiat) + "\n"
			myshare := new(big.Int).Div(new(big.Int).Mul(s.TotalAmount, big.NewInt(100)), total)
			info += fmt.Sprintf(doc.T("my_shares")+": %%%f.2\n", myshare)

			info += user.Address[0:6] + "..." + user.Address[len(user.Address)-4:]
			t.AddRow(info, doc.FormatValue(s.TotalAmount, 18), doc.FormatFiat(fiat))
		}
	}

	doc.WriteTable(t)

}
