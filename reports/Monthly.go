package reports

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlexNa-Holdings/savva-reports/assets"
	"github.com/AlexNa-Holdings/savva-reports/cmn"
	"github.com/AlexNa-Holdings/savva-reports/data"
	"github.com/AlexNa-Holdings/savva-reports/i18n"
	"github.com/AlexNa-Holdings/savva-reports/pdf"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

func Build(user_addr string, year, month int, output_path string, locale string) error {

	if month < 1 || month > 12 {
		log.Printf("Invalid month: %d", month)
		return fmt.Errorf("invalid month: %d", month)
	}

	doc, err := pdf.NewDoc(user_addr, locale)
	if err != nil {
		log.Printf("Error initializing PDF: %v", err)
		return fmt.Errorf("failed to initialize PDF: %w", err)
	}

	err = monthlyCoverPage(doc, year, month)
	if err != nil {
		log.Printf("Error creating cover page: %v", err)
		return fmt.Errorf("error creating cover page: %w", err)
	}

	addSectionLegal(doc)

	time_from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	time_to := time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.UTC)

	addSectionMyAuthors(doc, time_from, time_to)
	// addSectionSummary(doc, time_from, time_to)

	err = doc.WritePdf(output_path)
	if err != nil {
		log.Printf("Error writing PDF: %v", err)
		return fmt.Errorf("error saving PDF to %s: %w", output_path, err)
	}

	return nil
}

func monthlyCoverPage(doc *pdf.Doc, year, month int) error {
	doc.AddPage()

	user, err := data.GetUser(doc.UserAddress)
	if err != nil {
		log.Printf("Error fetching user data: %v", err)
		return fmt.Errorf("error fetching user data: %w", err)
	}

	// Draw avatar
	if err := doc.ImageFrom(user.AvatarImg, 60, 155, &gopdf.Rect{W: 500, H: 500}); err != nil {
		log.Error().Err(err).Msg("Failed to draw avatar image")
		return err
	}

	// print cmn.CoverImg on all page
	if err := doc.ImageFrom(assets.CoverImg, 0, 0, &gopdf.Rect{W: cmn.PageWidth, H: cmn.PageHeight}); err != nil {
		log.Error().Err(err).Msg("Failed to draw cover image")
		return err
	}

	//	doc.SetTextColor(0xff, 0x71, 0) //SAVVA
	doc.SetTextColor(0xff, 0xff, 0xff) //SAVVA
	doc.SetFont("DejaVuBold", "", 60)
	doc.TextCentered(fmt.Sprintf("%d", year), cmn.PageWidth-120, 70)
	doc.SetFont("DejaVuBold", "", 30)
	doc.TextCentered(strings.ToLower(i18n.GetMonthName(month, doc.Locale)), cmn.PageWidth-120, 100)

	doc.SetFont("DejaVuBold", "", 40)
	// doc.SetTextColor(0xc4, 0x58, 0) //Dark SAVVA
	doc.SetTextColor(0, 0, 0) //Dark SAVVA
	doc.TextCentered(strings.ToUpper(user.Name), 0, cmn.PageHeight-120)

	doc.SetFont("DejaVuBold", "", 20)
	// print address in form 0x1234...1234
	doc.TextCentered(user.Address[0:6]+"..."+user.Address[len(user.Address)-4:],
		cmn.PageWidth/2, cmn.PageHeight-95)

	doc.SetFont("Mono", "", 10)
	doc.TextCentered(fmt.Sprintf("Generated on: %s",
		time.Now().UTC().Format(time.RFC822)),
		0,
		cmn.PageHeight-70)

	doc.AddBlankPage()

	return nil
}
