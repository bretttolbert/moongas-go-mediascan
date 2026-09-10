package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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

	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Println("Error: Invalid arguments")
		fmt.Println("Usage: go run cmd/scanfilesyaml/main.go {config yaml filepath} {output yaml filepath} [mediaRootdir]")
		fmt.Println("Example: go run cmd/scanfilesyaml/main.go mediascan-config.yaml files.yaml /path/to/root")
		os.Exit(1)
	}
	configYamlFilepath := os.Args[1]
	outputYamlFilepath := os.Args[2]
	mediaRootDir := ""
	if len(os.Args) == 4 {
		mediaRootDir = os.Args[3]
	}

	conf := shared.LoadConf(configYamlFilepath)
	resolvedMediaDirs, err := shared.ResolveMediaDirs(mediaRootDir, conf.MediaDirs)
	if err != nil {
		log.Printf("ERROR: %v", err)
		os.Exit(1)
	}
	conf.MediaDirs = resolvedMediaDirs

	outputYamlAbsFilepath, err := filepath.Abs(outputYamlFilepath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to output file: %v", err)
	}

	var files shared.MediaFiles = shared.ScanFiles(conf)

	// consider "GroupBy" deprecated, I'm probably going to remove this feature
	// (GroupBy and MediaFilePlaylistList) as it duplicates functionality provided by mediaserver
	// however I was going to use for my timebox project, we shall see
	if conf.GroupBy == "year" {
		var mediaFilePlaylistList shared.MediaFilePlaylistList
		mediaFilePlaylistList.Playlists = make(map[string][]shared.MediaFile)
		for _, m := range files.Files {
			year := strconv.FormatInt(int64(m.Year), 10)
			_, ok := mediaFilePlaylistList.Playlists[year]
			if !ok {
				mediaFilePlaylistList.Playlists[year] = make([]shared.MediaFile, 0)
			}
			mediaFilePlaylistList.Playlists[year] = append(mediaFilePlaylistList.Playlists[year], m)
		}
		yamlData, err := yaml.Marshal(&mediaFilePlaylistList)
		shared.Check(err, "")
		err2 := os.WriteFile(outputYamlFilepath, yamlData, 0644)
		shared.Check(err2, outputYamlFilepath)
	} else {
		yamlData, err := yaml.Marshal(&files)
		shared.Check(err, "")
		err2 := os.WriteFile(outputYamlFilepath, yamlData, 0644)
		shared.Check(err2, outputYamlFilepath)
	}
	log.Printf("Written to file %s", outputYamlAbsFilepath)
}
