// Package render-doc implements documentation rendering for tasks.
// It is intended to be run via `go run`, passing a task / step action YAML manifest
// and a description in Asciidoctor format. The combined result will be
// written to the specified destination.
//
// Example invocation:
//
//	go run github.com/opendevstack/ods-pipeline/cmd/render-doc \
//		-manifest=tasks/my-task.yaml \
//		-description=build/docs/my-task.adoc \
//		-destination=docs/my-task.adoc
//
// By default, render-doc will use the template located at
// docs/template/tasks.adoc.tmpl to produce the resulting file. Another
// template can be specified via -template.
// To use the built-in template for step actions, pass -template=step-action.
// To use a custom template, pass -template=/path/to/my-custom-template.adoc.tmpl.
// If using a relative path, it has to be prefixed with "./".
package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/opendevstack/ods-pipeline/internal/projectpath"
	"github.com/opendevstack/ods-pipeline/pkg/renderdoc"
)

func main() {
	manifestFile := flag.String("manifest", "", "Manifest of Task or StepAction")
	descriptionFile := flag.String("description", "", "Description snippet")
	templateFile := flag.String("template", "task", "Template file. Either reference the built-in 'task' or 'step-action' templates, or specify a custom template file path.")
	destinationFile := flag.String("destination", "", "Destination file")
	flag.Parse()
	if err := render(*manifestFile, *descriptionFile, *templateFile, *destinationFile); err != nil {
		log.Fatal(err)
	}
}

func render(taskFile, descriptionFile, templateFile, destinationFile string) error {
	t, err := os.ReadFile(taskFile)
	if err != nil {
		return err
	}
	var description []byte
	if descriptionFile != "" {
		d, err := os.ReadFile(descriptionFile)
		if err != nil {
			return err
		}
		description = d
	}
	if !strings.HasPrefix(templateFile, "/") && !strings.HasPrefix(templateFile, "./") {
		templateFile = projectpath.RootedPath("docs/templates/" + templateFile + ".adoc.tmpl")
	}
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	task, err := renderdoc.ParseTask(t, description)
	if err != nil {
		return err
	}

	w, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	return renderdoc.RenderTaskDocumentation(w, tmpl, task)
}
