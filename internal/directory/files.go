package directory

import (
	"fmt"
	"io/ioutil"
	"log"
)

func ListFiles(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Files in %s:\n", dir)
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
