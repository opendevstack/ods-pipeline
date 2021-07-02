package config

import (
	"fmt"
	"io/ioutil"

	"sigs.k8s.io/yaml"
)

func GetODSConfig(filename string) (*ODS, error) {

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
