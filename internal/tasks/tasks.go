package tasks

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/command"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"sigs.k8s.io/yaml"
)

type Task struct {
	Name string
}

func parseTasks(helmTemplateOutput []byte) (map[string][]byte, error) {
	tasks := make(map[string][]byte)

	tasksBytes := bytes.Split(helmTemplateOutput, []byte("---"))

	for _, taskBytes := range tasksBytes {
		var t tekton.Task
		err := yaml.Unmarshal(taskBytes, &t)
		if err != nil {
			return nil, err
		}
		if len(t.Name) > 0 {
			tasks[t.Name] = taskBytes
		}
	}

	return tasks, nil
}

// Render extracts the task information into a struct, and
// executes the Asciidoctor template belonging to it.
func Render(sourceDir, targetDir string) error {
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return err
	}
	stdout, stderr, err := command.RunBufferedInDir(
		"helm",
		[]string{"template", "--values=values.docs.yaml", "."},
		sourceDir,
	)
	if err != nil {
		fmt.Println(string(stderr))
		return err
	}

	tasks, err := parseTasks(stdout)
	if err != nil {
		return err
	}
	for name, t := range tasks {
		targetFilename := fmt.Sprintf("%s.yaml", name)
		target := filepath.Join(targetDir, targetFilename)
		err := os.WriteFile(target, t, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
