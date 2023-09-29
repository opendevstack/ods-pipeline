// Package taskdoc implements documentation rendering for tasks.
// It is intended to be run via `go run`, passing a task YAML manifest
// and a description in Asciidoctor format. The combined result will be
// written to the specified destination.
//
// Example invocation:
//
//	go run github.com/opendevstack/ods-pipeline/cmd/taskdoc \
//		-task tasks/my-task.yaml \
//		-description build/docs/my-task.adoc \
//		-destination docs/my-task.adoc
//
// By default, taskdoc will use the template located at
// docs/tasks/template.adoc.tmpl to produce the resulting file. Another
// template can be specified via -template:
//
//	go run github.com/opendevstack/ods-pipeline/cmd/taskdoc \
//		-task tasks/my-task.yaml \
//		-description build/docs/my-task.adoc \
//		-template /path/to/my-custom-template.adoc.tmpl \
//		-destination docs/my-task.adoc
package main

import (
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/pkg/taskdoc"
)

func main() {
	taskFile := flag.String("task", "", "Task manifest")
	descriptionFile := flag.String("description", "", "Description snippet")
	templateFile := flag.String("template", projectpath.RootedPath("docs/tasks/template.adoc.tmpl"), "Template file")
	destinationFile := flag.String("destination", "", "Destination file")
	flag.Parse()
	if err := render(*taskFile, *descriptionFile, *templateFile, *destinationFile); err != nil {
		log.Fatal(err)
	}
}

func render(taskFile, descriptionFile, templateFile, destinationFile string) error {
	t, err := os.ReadFile(taskFile)
	if err != nil {
		return err
	}
	d, err := os.ReadFile(descriptionFile)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	task, err := taskdoc.ParseTask(t, d)
	if err != nil {
		return err
	}

	w, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	return taskdoc.RenderTaskDocumentation(w, tmpl, task)
}
