package taskdoc

import (
	"errors"
	"io"
	"text/template"

	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

type Param struct {
	Name        string
	Default     string
	Description string
}

type Result struct {
	Name        string
	Description string
}

type Task struct {
	Name        string
	Description string
	Params      []Param
	Results     []Result
}

// ParseTask reads a Tekton task from given bytes f,
// and assembles a new Task with the name, params and
// results from the parsed Tekton task, as well as the
// given description.
func ParseTask(f []byte, desc []byte) (*Task, error) {
	var t tekton.Task
	err := yaml.Unmarshal(f, &t)
	if err != nil {
		return nil, err
	}
	if t.Name == "" {
		return nil, errors.New("encountered empty name, something is wrong with the task")
	}
	task := &Task{
		Name:        t.Name,
		Description: string(desc),
		Params:      []Param{},
		Results:     []Result{},
	}
	for _, p := range t.Spec.Params {
		defaultValue := ""
		if p.Default != nil {
			defaultValue = p.Default.StringVal
		}
		task.Params = append(task.Params, Param{
			Name:        p.Name,
			Default:     defaultValue,
			Description: p.Description,
		})
	}
	for _, r := range t.Spec.Results {
		task.Results = append(task.Results, Result{
			Name:        r.Name,
			Description: r.Description,
		})
	}
	return task, nil
}

// RenderTaskDocumentation renders the given template with the task data,
// writing the result to w.
func RenderTaskDocumentation(w io.Writer, tmpl *template.Template, task *Task) error {
	if _, err := w.Write(
		[]byte("// File is generated; DO NOT EDIT.\n\n"),
	); err != nil {
		return err
	}
	return tmpl.Execute(w, task)
}
