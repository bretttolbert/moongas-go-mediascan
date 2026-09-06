package shared

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConf(configYamlFilepath string) (conf MediascanConf) {
	yfile, err := os.ReadFile(configYamlFilepath)
	Check(err, configYamlFilepath)
	err2 := yaml.Unmarshal(yfile, &conf)
	Check(err2, configYamlFilepath)
	return
}
