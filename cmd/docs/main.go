package main

import (
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/docs"
	"github.com/opendevstack/pipeline/internal/projectpath"
)

func main() {
	docs.RenderTasks(filepath.Join(projectpath.Root, "deploy/central/tasks"), filepath.Join(projectpath.Root, "docs/tasks"))
}
