package test

import (
	"context"
	"testing"
	"time"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"

	pipelineclientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getTr(ctx context.Context, t *testing.T, c pipelineclientset.Interface, name, ns string) (tr *v1beta1.TaskRun) {
	t.Helper()
	tr, err := c.TektonV1beta1().TaskRuns(ns).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		t.Error(err)
	}
	return tr
}

type conditionFn func(*v1beta1.TaskRun) bool

func waitForCondition(ctx context.Context, t *testing.T, c pipelineclientset.Interface, name, ns string, cond conditionFn, timeout time.Duration) *v1beta1.TaskRun {
	t.Helper()

	// Do a first quick check before setting the watch
	tr := getTr(ctx, t, c, name, ns)
	if cond(tr) {
		return tr
	}

	w, err := c.TektonV1beta1().TaskRuns(ns).Watch(ctx, metav1.SingleObject(metav1.ObjectMeta{
		Name:      name,
		Namespace: ns,
	}))
	if err != nil {
		t.Errorf("error watching taskrun: %s", err)
	}

	// Setup a timeout channel
	timeoutChan := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		timeoutChan <- struct{}{}
	}()

	// Wait for the condition to be true or a timeout
	for {
		select {
		case ev := <-w.ResultChan():
			tr := ev.Object.(*v1beta1.TaskRun)
			if cond(tr) {
				return tr
			}
		case <-timeoutChan:
			t.Fatal("time out")
		}
	}
}

func done(tr *v1beta1.TaskRun) bool {
	return tr.IsDone()
}
