package tektontaskrun

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	k "github.com/opendevstack/ods-pipeline/internal/kubernetes"
	"github.com/opendevstack/ods-pipeline/pkg/taskmanifest"
	tekton "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func installTask(path, namespace string, data map[string]string) (*tekton.Task, error) {
	var t tekton.Task
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("parse file: %w", err)
	}
	w := new(bytes.Buffer)
	err = taskmanifest.RenderTask(w, tmpl, data)
	if err != nil {
		return nil, fmt.Errorf("render task: %w", err)
	}
	err = yaml.Unmarshal(w.Bytes(), &t)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	clients := k.NewClients()
	tc := clients.TektonClientSet
	it, err := tc.TektonV1().Tasks(namespace).Create(context.TODO(), &t, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return it, nil
}
