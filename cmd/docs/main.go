package main

import (
	"log"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/docs"
	"github.com/opendevstack/pipeline/internal/projectpath"
)

func main() {
	err := docs.RenderTasks(
		filepath.Join(projectpath.Root, "deploy/ods-pipeline/charts/tasks"),
		filepath.Join(projectpath.Root, "docs/tasks"),
	)
	if err != nil {
		log.Fatal(err)
	}
}
