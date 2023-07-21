package main

import (
	"log"
	"path/filepath"

	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/internal/tasks"
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
