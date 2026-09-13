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

	if len(os.Args) < 3 || len(os.Args) > 4 {
		fmt.Println("Error: Invalid arguments")
		fmt.Println("Usage: go run cmd/scanfilesyaml/main.go {config yaml filepath} {output yaml filepath} [mediaRootdir]")
		fmt.Println("Example: go run cmd/scanfilesyaml/main.go mediascan-config.yaml files.yml /path/to/root")
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

	var files = shared.ScanFiles(conf)
	var filesYaml = shared.MediaFilesYamlFile{}
	filesYaml.Files = files
	yamlData, err := yaml.Marshal(&filesYaml)
	shared.Check(err, "")
	err2 := os.WriteFile(outputYamlFilepath, yamlData, 0644)
	shared.Check(err2, outputYamlFilepath)
	log.Printf("Written to file %s", outputYamlAbsFilepath)
}
