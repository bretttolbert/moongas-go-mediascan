// Package shared contains types shared across the mediascan application.
package shared

// ArtistYamlArtistData (DEPRECATED) includes the directory path and the data read
// from the artist.yml file in that directory.
// keep me in sync with python/mediascan/src/artist_yaml_artist_data.py
type ArtistYamlArtistData struct {
	ArtistData ArtistData `yaml:"artistData"`
	Path       string     `yaml:"path"`
}
