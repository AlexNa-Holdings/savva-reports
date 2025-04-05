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

func addSectionSponsored(doc *pdf.Doc, from, to time.Time) {

	if doc.History == nil {
		var err error
		doc.History, err = data.GetHistory(doc.UserAddress, &from, &to)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch history")
			return
		}
	}

	if doc.Sponsored == nil {
		var err error
		doc.Sponsored, err = data.GetSponsoredBy(doc.UserAddress)
		if err != nil {
			log.Error().Err(err).Msg("Failed to fetch sponsored by")
			return
		}
	}

	if len(doc.Sponsored) == 0 {
		log.Info().Msg("No sponsored data found")
		return // skip the section
	}

	total := big.NewInt(0)
	for _, s := range doc.Sponsored {
		total.Add(total, s.TotalAmount)
	}

	// Add a new section for the summary
	doc.NewSection(doc.T("sponsored.title"))

	doc.MarkDownToPdf(fmt.Sprintf(doc.T("sponsored.introduction"), doc.FormatValue(total, 18)) + " SAVVA (" + doc.FormatFiat(pdf.Value2Float(total, 18)*cmn.C.SavvaTokenPrice) + ")")
	doc.NewLine()

	t := pdf.NewTable()
	t.SetHeader(doc.T("account"), "SAVVA", cmn.C.CurrencySymbol)
	t.ColWidths = []float64{0, 100, 100}

	t.ColStyle[0].MinHeight = AVATAR_SIZE
	t.ColStyle[0].Padding.Left = AVATAR_SIZE + 10
	t.ColStyle[0].FontSize = 12
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

	for _, s := range doc.Sponsored {
		info := ""
		user, err := data.GetUser(s.Author)
		if err == nil {

			info += "!MD"
			if user.Name != "" {
				info += "## " + strings.ToUpper(user.Name) + "\u00AE\n"
			}

			dn := user.GetDisplayName()
			if dn != "" {
				info += dn + "\n"
			}

			info += user.Address[0:6] + "..." + user.Address[len(user.Address)-4:] + "\n"

			fiat_total_from_all := pdf.Value2Float(s.TotalFromAll, 18) * cmn.C.SavvaTokenPrice
			info += doc.T("total") + ": " + doc.FormatValue(s.TotalFromAll, 18) + " " + doc.FormatFiat(fiat_total_from_all) + "\n"

			if s.TotalFromAll != nil && s.TotalFromAll.Cmp(big.NewInt(0)) != 0 { // just to be sure
				myshare10 := new(big.Int).Div(new(big.Int).Mul(s.TotalAmount, big.NewInt(10000)), s.TotalFromAll)
				myshare := float64(myshare10.Int64()) / 100.0
				info += fmt.Sprintf(doc.T("my_share")+": %.2f%%\n", myshare)
			}
		}

		fiat := pdf.Value2Float(s.TotalAmount, 18) * cmn.C.SavvaTokenPrice
		t.AddRow(info, doc.FormatValue(s.TotalAmount, 18), doc.FormatFiat(fiat))

	}

	doc.WriteTable(t)
}
