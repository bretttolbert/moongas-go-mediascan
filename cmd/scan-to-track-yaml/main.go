package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bretttolbert/moongas-mediascan-go/internal/shared"
	"gopkg.in/yaml.v3"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Unable to determine current working directory: %v", err)
	} else {
		log.Printf("Current working directory: %s", cwd)
	}

	if len(os.Args) < 2 || len(os.Args) > 3 {
		log.Println("Error: Invalid arguments")
		log.Println("Usage: go run cmd/scan-to-track-yaml <config-yaml> [media-root]")
		log.Println("Example: go run cmd/scan-to-track-yaml mediascan-config.yml /path/to/root")
		os.Exit(1)
	}
	configYamlFilepath := os.Args[1]
	mediaRootDir := ""
	if len(os.Args) == 3 {
		mediaRootDir = os.Args[2]
	}

	conf := shared.LoadConf(configYamlFilepath)
	resolvedMediaDirs, err := shared.ResolveMediaDirs(mediaRootDir, conf.MediaDirs)
	if err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
	conf.MediaDirs = resolvedMediaDirs

	files := shared.ScanFiles(conf)

	countCreated := 0
	countSkipped := 0
	for _, m := range files {
		ext := filepath.Ext(m.Path)
		trackYamlFilepath := filepath.Join(m.AlbumPath, strings.TrimSuffix(filepath.Base(m.Path), ext)+".yml")
		if _, err := os.Stat(trackYamlFilepath); err == nil {
			log.Printf("Skipping %s (%s already exists)", m.Path, trackYamlFilepath)
			countSkipped += 1
			continue
		}

		trackYaml := shared.TrackYamlFile{
			TrackData: shared.TrackData{
				Title:       m.Title,
				Artist:      m.Artist,
				AlbumArtist: m.AlbumArtist,
				Album:       m.Album,
				Genre:       m.Genre,
				Year:        m.Year,
			},
		}
		yamlData, err := yaml.Marshal(&trackYaml)
		shared.Check(err, trackYamlFilepath)
		err = os.WriteFile(trackYamlFilepath, yamlData, 0644)
		shared.Check(err, trackYamlFilepath)
		log.Printf("Created %s", trackYamlFilepath)
		countCreated += 1
	}

	log.Printf("Scanned %d media files", len(files))
	log.Printf("Created %d track YAML (.yml) files", countCreated)
	log.Printf("Skipped %d track YAML (.yml) files (already exist)", countSkipped)
}
