package shared

// Data model for artists.yml file output by mediascan cmd/scanartists
// keep me in sync with python/mediascan/src/artists.py
type Artists struct {
	Artists []Artist
}
