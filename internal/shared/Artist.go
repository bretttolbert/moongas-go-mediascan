package shared

// an artist, including directory path and the data read from the artist.yaml file in said directory
// keep me in sync with python/mediascan/src/artist.py
type Artist struct {
	ArtistData ArtistData `yaml:"artistData"`
	Path       string     `yaml:"path"`
}
