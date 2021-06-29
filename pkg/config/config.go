package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func GetODSConfig(project, repository, gitFullRef, filename string) (*ODS, error) {

	odsFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var odsConfig ODS
	err = yaml.Unmarshal(odsFile, &odsConfig)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config: %w", err)
	}
	return &odsConfig, nil
}
