package sonar

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(baseURL string) *Client {
	return NewClient(&ClientConfig{BaseURL: baseURL})
}

func TestQualityGateGet(t *testing.T) {

	srv := newTestServer([]recordedResponse{})
	defer srv.Server.Close()

	c := testClient(srv.Server.URL)

	tests := map[string]struct {
		Fixture    string
		WantStatus string
	}{
		"ERROR status": {
			Fixture:    "project_status_error.json",
			WantStatus: "ERROR",
		},
		"OK status": {
			Fixture:    "project_status_ok.json",
			WantStatus: "OK",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := ioutil.ReadFile("../../test/testdata/sonar/fixtures/" + tc.Fixture)
			if err != nil {
				t.Fatal(err)
			}
			srv.addResponse(recordedResponse{Body: b, StatusCode: 200})
			got, err := c.QualityGateGet(QualityGateGetParams{})
			if err != nil {
				t.Fatalf("Unexpected error on request: %s", err)
			}
			lastRequest := srv.getLastRequest()
			if lastRequest.Path != "/api/qualitygates/project_status" {
				t.Fatalf("want req path %s, got %s", tc.WantStatus, lastRequest.Path)
			}
			if got.ProjectStatus.Status != tc.WantStatus {
				t.Fatalf("want %s, got %s", tc.WantStatus, got.ProjectStatus.Status)
			}
		})
	}
}

type recordedRequest struct {
	Body []byte
	Path string
}

type recordedResponse struct {
	Body       []byte
	StatusCode int
}

type testServer struct {
	Requests  []recordedRequest
	Responses []recordedResponse
	Server    *httptest.Server
}

func (ts *testServer) handler(w http.ResponseWriter, r *http.Request) {
	ts.Requests = append(ts.Requests, recordedRequest{
		Path: r.URL.Path,
	})
	var res recordedResponse
	res, ts.Responses = ts.Responses[0], ts.Responses[1:]
	w.WriteHeader(res.StatusCode)
	w.Write(res.Body)
}

func (ts *testServer) addResponse(res recordedResponse) {
	ts.Responses = append(ts.Responses, res)
}

func (ts *testServer) getLastRequest() recordedRequest {
	return ts.Requests[len(ts.Requests)-1]
}

func newTestServer(responses []recordedResponse) *testServer {
	ts := &testServer{
		Responses: responses,
	}
	ts.Server = httptest.NewServer(http.HandlerFunc(ts.handler))
	return ts
}
