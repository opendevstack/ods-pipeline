package sonar

import (
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestQualityGateGet(t *testing.T) {

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	c := testClient(srv.Server.URL)

	tests := map[string]struct {
		Fixture    string
		WantStatus string
	}{
		"ERROR status": {
			Fixture:    "sonar/project_status_error.json",
			WantStatus: "ERROR",
		},
		"OK status": {
			Fixture:    "sonar/project_status_ok.json",
			WantStatus: "OK",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv.EnqueueResponse(
				t, "/api/qualitygates/project_status",
				200, tc.Fixture,
			)
			got, err := c.QualityGateGet(QualityGateGetParams{})
			if err != nil {
				t.Fatalf("Unexpected error on request: %s", err)
			}
			if got.ProjectStatus.Status != tc.WantStatus {
				t.Fatalf("want %s, got %s", tc.WantStatus, got.ProjectStatus.Status)
			}
		})
	}
}
