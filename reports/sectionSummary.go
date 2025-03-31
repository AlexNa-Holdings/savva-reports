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

	var fiat float64

	c := calcCounters(doc)

	fiat = pdf.Value2Float(c.savva_in, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.savva_in"), doc.FormatValue(c.savva_in, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.savva_out, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.savva_out"), doc.FormatValue(c.savva_out, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.donations_contribute, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.donations_contribute"), doc.FormatValue(c.donations_contribute, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.donations_received, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.donations_received"), doc.FormatValue(c.donations_received, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.fund_contributed, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.fund_contributed"), doc.FormatValue(c.fund_contributed, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.fund_prizes_won, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.fund_prizes_won"), doc.FormatValue(c.fund_prizes_won, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.staking_in, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.staking_in"), doc.FormatValue(c.staking_in, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.staking_out, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.staking_out"), doc.FormatValue(c.staking_out, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.staking_staked, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.staking_staked"), doc.FormatValue(c.staking_staked, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.club_buy, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.club_buy"), doc.FormatValue(c.club_buy, 18), doc.FormatFiat(fiat))
	fiat = pdf.Value2Float(c.club_claimed, 18) * cmn.C.SavvaTokenPrice
	t.AddRow(doc.T("summary.club_claimed"), doc.FormatValue(c.club_claimed, 18), doc.FormatFiat(fiat))

	doc.WriteTable(t)

}

type Counters struct {
	savva_in             *big.Int
	savva_out            *big.Int
	donations_contribute *big.Int
	donations_received   *big.Int
	fund_contributed     *big.Int
	fund_prizes_won      *big.Int
	staking_in           *big.Int
	staking_out          *big.Int
	staking_staked       *big.Int
	club_buy             *big.Int
	club_claimed         *big.Int
}

func calcCounters(doc *pdf.Doc) *Counters {
	c := &Counters{
		savva_in:             new(big.Int),
		savva_out:            new(big.Int),
		donations_contribute: new(big.Int),
		donations_received:   new(big.Int),
		fund_contributed:     new(big.Int),
		fund_prizes_won:      new(big.Int),
		staking_in:           new(big.Int),
		staking_out:          new(big.Int),
		staking_staked:       new(big.Int),
		club_buy:             new(big.Int),
		club_claimed:         new(big.Int),
	}

	for _, h := range doc.History {

		amount := new(big.Int)
		if h.Amount != nil {
			amount = h.Amount
		}

		if h.Contract.String == "token" && h.Type.String == "transfer" {
			if h.FromAddr.String == doc.UserAddress {
				c.savva_out = c.savva_out.Sub(c.savva_out, amount)
			} else if h.ToAddr.String == doc.UserAddress {
				c.savva_in = c.savva_in.Add(c.savva_in, amount)
			}
		}

		if h.Contract.String == "fund" {
			if h.Type.String == "donation" {
				if h.FromAddr.String == doc.UserAddress {
					c.donations_contribute = c.donations_contribute.Sub(c.donations_contribute, amount)
				}

				if h.ToAddr.String == doc.UserAddress {
					c.donations_received = c.donations_received.Add(c.donations_received, amount)
				}
			}

			if h.Type.String == "contribute" {
				if h.FromAddr.String == doc.UserAddress {
					c.fund_contributed = c.fund_contributed.Sub(c.fund_contributed, amount)
				}
			}

			if h.Type.String == "prize" {
				if h.ToAddr.String == doc.UserAddress {
					c.fund_prizes_won = c.fund_prizes_won.Add(c.fund_prizes_won, amount)
				}
			}
		}

		if h.Contract.String == "staking" {
			if h.Type.String == "transferred" {
				if h.FromAddr.String == doc.UserAddress {
					c.staking_out = c.staking_out.Sub(c.staking_out, amount)
				} else if h.ToAddr.String == doc.UserAddress {
					c.staking_in = c.staking_in.Add(c.staking_in, amount)
				}
			}

			if h.Type.String == "staked" {
				if h.FromAddr.String == doc.UserAddress {
					c.staking_staked = c.staking_staked.Sub(c.staking_staked, amount)
				}
			}

			if h.Type.String == "us_claimed" {
				if h.ToAddr.String == doc.UserAddress {
					c.staking_staked = c.staking_staked.Add(c.staking_staked, amount)
				}
			}
		}

		if h.Contract.String == "club" {
			if h.Type.String == "buy" {
				if h.FromAddr.String == doc.UserAddress {
					c.club_buy = c.club_buy.Sub(c.club_buy, amount)
				}
			}

			if h.Type.String == "claimed" {
				if h.ToAddr.String == doc.UserAddress {
					c.club_claimed = c.club_claimed.Add(c.club_claimed, amount)
				}
			}
		}
	}

	return c
}
