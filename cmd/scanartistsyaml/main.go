package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bretttolbert/moongas-mediascan-go/internal/shared"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Error: Invalid arguments")
		fmt.Println("Usage: go run cmd/scanartistsyaml/main.go {config yaml filepath} {output yaml filepath}")
		fmt.Println("Example: go run cmd/scanartistsyaml/main.go conf/conf.yaml out/artists.yaml")
		os.Exit(1)
	}
	configYamlFilepath := os.Args[1]
	outputYamlFilepath := os.Args[2]
	conf := shared.LoadConf(configYamlFilepath)
	var artists shared.Artists = shared.ScanArtists(conf)

	yamlData, err := yaml.Marshal(&artists)
	shared.Check(err, "")
	err2 := os.WriteFile(outputYamlFilepath, yamlData, 0644)
	shared.Check(err2, outputYamlFilepath)

	log.Printf("Successfully loaded %d artist.yaml files", len(artists.Artists))
	log.Printf("Written to file %s", outputYamlFilepath)
}
