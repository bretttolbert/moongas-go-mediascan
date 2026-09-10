package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	if len(os.Args) != 3 {
		fmt.Println("Error: Invalid arguments")
		fmt.Println("Usage: go run cmd/scanartistsyaml/main.go {config yaml filepath} {output yaml filepath}")
		fmt.Println("Example: go run cmd/scanartistsyaml/main.go mediascan-config.yaml artists.yaml")
		os.Exit(1)
	}
	configYamlFilepath := os.Args[1]
	outputYamlFilepath := os.Args[2]

	conf := shared.LoadConf(configYamlFilepath)

	outputYamlAbsFilepath, err := filepath.Abs(outputYamlFilepath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to output file: %v", err)
	}

	var artists shared.Artists = shared.ScanArtists(conf)

	yamlData, err := yaml.Marshal(&artists)
	shared.Check(err, "")
	err2 := os.WriteFile(outputYamlAbsFilepath, yamlData, 0644)
	shared.Check(err2, outputYamlAbsFilepath)

	log.Printf("Successfully loaded %d artist.yaml files", len(artists.Artists))
	log.Printf("Written to file %s", outputYamlAbsFilepath)
}
