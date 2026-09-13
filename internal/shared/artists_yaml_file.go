package shared

// ArtistsYamlFile is the data model for the artists.yml file output by mediascan cmd/scanartists.
// keep me in sync with python/mediascan/src/artists.py
type ArtistsYamlFile struct {
	Artists []ArtistYamlArtistData
}
