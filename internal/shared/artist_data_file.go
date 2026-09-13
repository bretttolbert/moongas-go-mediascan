package shared

// ArtistDataFile represents an artist.yml file.
// keep me in sync with python/mediascan/src/artistdatafile.py
type ArtistDataFile struct {
	ArtistData ArtistData `yaml:"artistData"`
}
