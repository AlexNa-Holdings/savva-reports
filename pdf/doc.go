package pdf

import (
	"github.com/AlexNa-Holdings/savva-reports/assets"
	"github.com/AlexNa-Holdings/savva-reports/i18n"
	"github.com/rs/zerolog/log"
	"github.com/signintech/gopdf"
)

type Doc struct { // Extended gopdf.GoPdf
	*gopdf.GoPdf
	Locale      string
	CurentPage  int
	UserAddress string
}

func NewDoc(user_addr, locale string) (*Doc, error) {
	// Create a new PDF document.
	doc := Doc{
		GoPdf:       new(gopdf.GoPdf),
		UserAddress: user_addr,
		Locale:      locale,
		CurentPage:  0,
	}
	doc.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// Load all fonts
	for name, font := range assets.AllFonts {
		err := doc.AddTTFFontData(name, font)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to load font %s", name)
			return nil, err
		}
	}

	// Set default font.  Handle the error.
	if err := doc.SetFont("Arial", "", 24); err != nil {
		log.Error().Err(err).Msg("Failed to set font for cover page")
		return nil, err
	}
	doc.SetTextColor(0, 0, 0) // Black

	return &doc, nil
}

func (doc *Doc) T(key string) string {
	return i18n.T(key, doc.Locale)
}
