package sonar

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestComputeEngineTaskGet(t *testing.T) {

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	c := testClient(srv.Server.URL)

	tests := map[string]struct {
		Fixture    string
		WantStatus string
	}{
		"FAILED status": {
			Fixture:    "sonar/task_failed.json",
			WantStatus: TaskStatusFailed,
		},
		"SUCCESS status": {
			Fixture:    "sonar/task_success.json",
			WantStatus: TaskStatusSuccess,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv.EnqueueResponse(
				t, "/api/ce/task",
				200, tc.Fixture,
			)
			taskID := "AVAn5RKqYwETbXvgas-I"
			got, err := c.ComputeEngineTaskGet(ComputeEngineTaskGetParams{ID: taskID})
			if err != nil {
				t.Fatalf("Unexpected error on request: %s", err)
			}

			// check extracted status matches
			if got.Status != tc.WantStatus {
				t.Fatalf("want %s, got %s", tc.WantStatus, got.Status)
			}

			// check sent task ID matches
			lr, err := srv.LastRequest()
			if err != nil {
				t.Fatal(err)
			}
			q := lr.URL.Query()
			if q.Get("id") != taskID {
				t.Fatalf("want %s, got %s", taskID, q.Get("id"))
			}
		})
	}
}
