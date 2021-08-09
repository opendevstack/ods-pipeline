package testserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/logging"
)

type RecordedResponse struct {
	Body       []byte
	StatusCode int
	Fixture    string
}

type TestServer struct {
	Responses map[string][]RecordedResponse
	Server    *httptest.Server
}

func (ts *TestServer) stub(l logging.SimpleLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if res, ok := ts.Responses[r.URL.Path]; ok {
			if len(res) > 0 {
				response := res[0]
				l.Logf("responding with body from %s", response.Fixture)
				ts.Responses[r.URL.Path] = res[1:]
				w.WriteHeader(response.StatusCode)
				_, err := w.Write(response.Body)
				if err != nil {
					http.Error(w, "write error", http.StatusInternalServerError)
					return
				}
				return
			}
		}
		l.Logf("no stub response registered for path %s", r.URL.Path)
		http.NotFound(w, r)
	}
}

func (ts *TestServer) EnqueueResponse(t *testing.T, path string, statusCode int, fixture string) {
	body := []byte{}
	if len(fixture) > 0 {
		filename := filepath.Join(projectpath.Root, "test/testdata/fixtures", fixture)
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}
		body = b
	}
	if _, ok := ts.Responses[path]; !ok {
		ts.Responses[path] = []RecordedResponse{}
	}
	ts.Responses[path] = append(ts.Responses[path], RecordedResponse{
		Body:       body,
		StatusCode: statusCode,
		Fixture:    fixture,
	})
}

func NewTestServer(l logging.SimpleLogger) (*TestServer, func()) {
	ts := &TestServer{
		Responses: make(map[string][]RecordedResponse),
	}
	ts.Server = httptest.NewServer(http.HandlerFunc(ts.stub(l)))
	return ts, func() { ts.Server.Close() }
}
