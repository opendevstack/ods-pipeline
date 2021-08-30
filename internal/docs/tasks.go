package docs

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/opendevstack/pipeline/internal/command"
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
	templateFilename := filepath.Join(targetDir, "template.adoc.tmpl")
	templateFileParts := strings.Split(templateFilename, "/")
	templateDisplayname := templateFileParts[len(templateFileParts)-1]
	_, err = targetFile.WriteString(
		"// Document generated by internal/documentation/tasks.go from " + templateDisplayname + "; DO NOT EDIT.\n\n",
	)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(templateFilename)
	if err != nil {
		return err
	}
	return tmpl.Execute(targetFile, data)
}

func parseTasks(helmTemplateOutput []byte) ([]*tekton.ClusterTask, error) {
	var tasks []*tekton.ClusterTask

	tasksBytes := bytes.Split(helmTemplateOutput, []byte("---"))

	for _, taskBytes := range tasksBytes {
		var t tekton.ClusterTask
		err := yaml.Unmarshal(taskBytes, &t)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal tasks: %w", err)
		}
		if len(t.Name) > 0 {
			tasks = append(tasks, &t)
		}
	}

	return tasks, nil
}

// RenderTasks extracts the task information into a struct, and
// executes the Asciidoctor template belonging to it.
func RenderTasks(sourceDir, targetDir string) error {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return err
	}
	stdout, stderr, err := command.RunInDir(
		"helm",
		[]string{"template", "--values=values.docs.yaml", "."},
		sourceDir,
	)
	if err != nil {
		fmt.Println(string(stderr))
		log.Fatal(err)
	}

	tasks, err := parseTasks(stdout)
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tasks {
		task := Task{
			Name:        t.Name,
			Description: t.Spec.Description,
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
		err := renderTemplate(targetDir, target, task)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
