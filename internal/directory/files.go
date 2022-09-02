package directory

import (
	"fmt"
	"log"
	"os"
)

func ListFiles(dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Files in %s:\n", dir)
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
