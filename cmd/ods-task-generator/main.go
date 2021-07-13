package main

import (
	"fmt"
	"os"

	"github.com/opendevstack/pipeline/internal/odstaskgen"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Please provide a technology")
		return
	}

	argsWithoutProg := os.Args[1:]
	technology := argsWithoutProg[0]

	odstaskgen.GenerateODSTaskFiles(technology)
}
