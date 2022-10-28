package main

import (
	"fmt"
	"os"
)

const (
	// ageKeyFilePath is the path where to store the age-key-secret openshift secret content,
	// required by the helm secrets plugin.
	ageKeyFilePath = "./key.txt"
)

func storeAgeKey(ageKeyContent []byte) error {
	file, err := os.Create(ageKeyFilePath)
	if err != nil {
		return fmt.Errorf("create age key file path: %w", err)
	}
	defer file.Close()
	_, err = file.Write(ageKeyContent)
	if err != nil {
		return fmt.Errorf("write age key: %w", err)
	}
	return err
}
