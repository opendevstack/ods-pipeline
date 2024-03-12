// Package render-manifest implements manifest rendering for K8s manifests.
// It is intended to be run via `go run`, passing a YAML template
// and data to be rendered. The combined result will be
// written to the specified destination. The -data flag can be passed
// multiple times and may specify any key-value combination, which can then
// be consumed in the template through Go's text/template package. E.g.
// passing -data Foo=bar will replace {{.Foo}} in the template with bar.
//
// Example invocation:
//
//	go run github.com/opendevstack/ods-pipeline/cmd/render-manifest \
//		-data ImageRepository=ghcr.io/my-org/my-repo \
//		-data Version=latest \
//		-template build/tasks/my-task.yaml \
//		-destination tasks/my-task.yaml
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/opendevstack/ods-pipeline/pkg/rendermanifest"
	"github.com/opendevstack/ods-pipeline/pkg/tektontaskrun"
)

func main() {
	templateFile := flag.String("template", "", "Template file")
	destinationFile := flag.String("destination", "", "Destination file")
	cc := tektontaskrun.NewClusterConfig()
	mf := &MapFlag{v: cc.DefaultManifestTemplateData()}
	flag.Var(mf, "data", "Key-value pairs")
	flag.Parse()
	if err := render(*templateFile, *destinationFile, mf.v); err != nil {
		log.Fatal(err)
	}
}

func render(templateFile, destinationFile string, data map[string]string) error {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	w, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	return rendermanifest.RenderManifest(w, tmpl, data)
}

type MapFlag struct {
	v map[string]string
}

func (mf *MapFlag) String() string {
	return fmt.Sprintf("%v", mf.v)
}
func (mf *MapFlag) Set(v string) error {
	key, value, ok := strings.Cut(v, "=")
	if !ok {
		return fmt.Errorf("must have = sign")
	}
	mf.v[key] = value
	return nil
}
