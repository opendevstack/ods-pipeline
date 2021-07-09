package file

import (
	"io"
	"os"
)

// Copy a file
func Copy(from, to string) error {
	// Open original file
	original, err := os.Open(from)
	if err != nil {
		return err
	}
	defer original.Close()

	// Create new file
	new, err := os.Create(to)
	if err != nil {
		return err
	}
	defer new.Close()

	// Copy contents of original file into new file
	_, err = io.Copy(new, original)
	if err != nil {
		return err
	}
	return nil
}
