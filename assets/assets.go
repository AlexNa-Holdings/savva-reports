package assets

import (
	"bytes"
	_ "embed" // for embedding assets
	"image"

	"github.com/rs/zerolog/log"
)

// fonts
//
//go:embed fonts/Arial.ttf
var FontAreal []byte

//go:embed fonts/Times.ttf
var FontTimes []byte

//go:embed fonts/Mono.ttf
var FontMono []byte

//go:embed fonts/DejaVuBold.ttf
var FontDejaVuBold []byte

var AllFonts = map[string][]byte{
	"Arial":      FontAreal,
	"Times":      FontTimes,
	"Mono":       FontMono,
	"DejaVuBold": FontDejaVuBold,
}

// images
//
//go:embed images/avatar-default.png
var AvatarDefault []byte
var AvatarDefaultImg image.Image

//go:embed images/SAVVA.png
var LogoSavva []byte
var LogoSavvaImg image.Image

//go:embed images/cover.png
var Cover []byte
var CoverImg image.Image

func ProcessAssets() error {
	var err error

	LogoSavvaImg, _, err = image.Decode(bytes.NewReader(LogoSavva))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode logo image for background")
		return err
	}

	CoverImg, _, err = image.Decode(bytes.NewReader(Cover))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode cover image for background")
		return err
	}

	AvatarDefaultImg, _, err = image.Decode(bytes.NewReader(AvatarDefault))
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode default avatar image")
		return err
	}

	return nil
}

// init the package
func init() {
	// Process assets
	if err := ProcessAssets(); err != nil {
		log.Fatal().Err(err).Msg("Failed to process assets")
	}
}
