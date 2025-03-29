package cmn

import (
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Load file from the IPFS gateway
func Ipfs(cid string) []byte {

	url := C.IPFS_Gateway + cid

	resp, err := http.Get(url)
	if err != nil {
		log.Error().Err(err).Msgf("Error fetching IPFS file: %v", err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Error fetching IPFS file: %s", resp.Status)
		return nil
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading IPFS file")
		return nil
	}
	return data
}
