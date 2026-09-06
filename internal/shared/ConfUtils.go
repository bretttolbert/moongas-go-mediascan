package shared

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func LoadConf(configYamlFilepath string) (conf MediascanConf) {
	configYamlAbsFilepath, err := filepath.Abs(configYamlFilepath)
	if err != nil {
		log.Fatalf("Failed to get absolute path to yaml config file: %v", err)
	} else {
		log.Printf("Loading config from yaml file: %v", configYamlAbsFilepath)
	}
	yfile, err := os.ReadFile(configYamlAbsFilepath)
	Check(err, configYamlAbsFilepath)
	err2 := yaml.Unmarshal(yfile, &conf)
	Check(err2, configYamlAbsFilepath)
	return
}
