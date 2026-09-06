package shared

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func ScanArtists(conf MediascanConf) Artists {
	var artists Artists
	countLoadFailed := 0
	countSkipped := 0
	for _, mediaDir := range conf.MediaDirs {
		err := filepath.Walk(mediaDir,
			func(path string, info os.FileInfo, err error) error {

				if !info.IsDir() && info.Name() == "artist.yaml" {

					log.Printf("Reading filepath %s", path)
					if err != nil {
						return err
					}
					if ContainsAnyOf(path, conf.ExcludePaths) {
						log.Printf("Skipping %s (ExcludePaths)", path)
						countSkipped += 1
						return nil
					}
					var a Artist
					a.Path = filepath.Dir(path)

					// Read the YAML file
					data, err := os.ReadFile(path)
					Check(err, path)

					// Unmarshal the YAML data into the Config struct
					var dataFile ArtistDataFile
					err = yaml.Unmarshal(data, &dataFile)
					Check(err, path)

					a.ArtistData = dataFile.ArtistData

					artists.Artists = append(artists.Artists, a)
				}
				return nil
			})
		if err != nil {
			log.Printf("ERROR loading file: %v", err)
			countLoadFailed += 1
		}
	}
	return artists
}
