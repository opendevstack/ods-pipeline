package docs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/opendevstack/ods-pipeline/internal/command"
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

func renderTemplate(targetDir, targetFilename string, data Task) error {
	targetFile, err := os.Create(targetFilename)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(filepath.Join(targetDir, "template.adoc.tmpl"))
	if err != nil {
		return err
	}
	return RenderTaskDocumentation(targetFile, tmpl, &data)
}

func parseTasks(helmTemplateOutput []byte) ([]*tekton.Task, error) {
	var tasks []*tekton.Task

	tasksBytes := bytes.Split(helmTemplateOutput, []byte("---"))

	for _, taskBytes := range tasksBytes {
		var t tekton.Task
		err := yaml.Unmarshal(taskBytes, &t)
		if err != nil {
			return nil, err
		}
		if len(t.Name) > 0 {
			tasks = append(tasks, &t)
		}
	}

	return tasks, nil
}

// RenderTasks extracts the task information into a struct, and
// executes the Asciidoctor template belonging to it.
func RenderTasks(tasksSourceDir, descriptionsSourceDir, targetDir string) error {
	if _, err := os.Stat(tasksSourceDir); os.IsNotExist(err) {
		return err
	}
	if _, err := os.Stat(descriptionsSourceDir); os.IsNotExist(err) {
		return err
	}
	stdout, stderr, err := command.RunBufferedInDir(
		"helm",
		[]string{"template", "--values=values.docs.yaml", "."},
		tasksSourceDir,
	)
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}

	tasks, err := parseTasks(stdout)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		desc, err := os.ReadFile(filepath.Join(descriptionsSourceDir, fmt.Sprintf("%s.adoc", t.Name)))
		if err != nil {
			return err
		}
		task := Task{
			Name:        t.Name,
			Description: string(desc),
			Params:      []Param{},
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
		targetFilename := fmt.Sprintf("%s.adoc", t.Name)
		target := filepath.Join(targetDir, targetFilename)
		err = renderTemplate(targetDir, target, task)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

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

func RenderTaskDocumentation(w io.Writer, tmpl *template.Template, task *Task) error {
	if _, err := w.Write(
		[]byte("// File is generated; DO NOT EDIT.\n\n"),
	); err != nil {
		return err
	}
	return tmpl.Execute(w, task)
}
