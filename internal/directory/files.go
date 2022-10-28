package directory

import (
	"fmt"
	"io"
	"os"
)

func ListFiles(dir string, out io.Writer) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Files in %s:\n", dir)
	for _, file := range files {
		fmt.Fprintln(out, file.Name())
	}
	return nil
}
