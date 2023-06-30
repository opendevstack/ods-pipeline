package sonar

import (
	"testing"

	"github.com/opendevstack/ods-pipeline/test/testserver"
)

func TestQualityGateGet(t *testing.T) {

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	c := testClient(t, srv.Server.URL)

	tests := map[string]struct {
		responseFixture string
		params          QualityGateGetParams
		wantRequestURI  string
		wantStatus      string
	}{
		"ERROR status": {
			params:          QualityGateGetParams{ProjectKey: "foo"},
			responseFixture: "sonar/project_status_error.json",
			wantRequestURI:  "/api/qualitygates/project_status?projectKey=foo",
			wantStatus:      "ERROR",
		},
		"OK status": {
			params:          QualityGateGetParams{ProjectKey: "foo"},
			responseFixture: "sonar/project_status_ok.json",
			wantRequestURI:  "/api/qualitygates/project_status?projectKey=foo",
			wantStatus:      "OK",
		},
		"OK status for branch": {
			params:          QualityGateGetParams{ProjectKey: "foo", Branch: "bar"},
			responseFixture: "sonar/project_status_ok.json",
			wantRequestURI:  "/api/qualitygates/project_status?projectKey=foo&branch=bar",
			wantStatus:      "OK",
		},
		"OK status for branch (PR=0)": {
			params:          QualityGateGetParams{ProjectKey: "foo", Branch: "bar", PullRequest: "0"},
			responseFixture: "sonar/project_status_ok.json",
			wantRequestURI:  "/api/qualitygates/project_status?projectKey=foo&branch=bar",
			wantStatus:      "OK",
		},
		"OK status for PR": {
			params:          QualityGateGetParams{ProjectKey: "foo", PullRequest: "123"},
			responseFixture: "sonar/project_status_ok.json",
			wantRequestURI:  "/api/qualitygates/project_status?projectKey=foo&pullRequest=123",
			wantStatus:      "OK",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv.EnqueueResponse(
				t, "/api/qualitygates/project_status",
				200, tc.responseFixture,
			)
			got, err := c.QualityGateGet(tc.params)
			if err != nil {
				t.Fatalf("Unexpected error on request: %s", err)
			}
			if got.ProjectStatus.Status != tc.wantStatus {
				t.Fatalf("want %s, got %s", tc.wantStatus, got.ProjectStatus.Status)
			}
			req, err := srv.LastRequest()
			if err != nil {
				t.Fatal(err)
			}
			if req.URL.RequestURI() != tc.wantRequestURI {
				t.Fatalf("want request URI %s, got %s", tc.wantRequestURI, req.URL.RequestURI())
			}
		})
	}
}
