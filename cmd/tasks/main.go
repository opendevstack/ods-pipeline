package main

import (
	"log"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/internal/tasks"
)

func main() {
	err := tasks.Render(
		filepath.Join(projectpath.Root, "deploy/ods-pipeline/charts/tasks"),
		filepath.Join(projectpath.Root, "tasks"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
