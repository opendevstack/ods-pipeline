package main

import (
	"log"
	"path/filepath"

	"github.com/opendevstack/ods-pipeline/internal/docs"
	"github.com/opendevstack/ods-pipeline/internal/projectpath"
)

func main() {
	err := docs.RenderTasks(
		filepath.Join(projectpath.Root, "deploy/ods-pipeline/charts/tasks"),
		filepath.Join(projectpath.Root, "docs/tasks/descriptions"),
		filepath.Join(projectpath.Root, "docs/tasks"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
