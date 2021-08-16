package testserver

import (
	"bytes"
	"errors"
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
	EnqueuedResponses map[string][]RecordedResponse
	ObservedRequests  []*http.Request
	Server            *httptest.Server
}

func cloneRequest(r *http.Request) (*http.Request, error) {
	r2 := *r
	var b bytes.Buffer
	_, err := b.ReadFrom(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(&b)
	r2.Body = ioutil.NopCloser(bytes.NewReader(b.Bytes()))
	return &r2, nil
}

func (ts *TestServer) stub(l logging.SimpleLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cr, err := cloneRequest(r)
		if err != nil {
			http.Error(w, "clone request error", http.StatusInternalServerError)
			return
		}
		ts.ObservedRequests = append(ts.ObservedRequests, cr)
		if res, ok := ts.EnqueuedResponses[r.URL.Path]; ok {
			if len(res) > 0 {
				response := res[0]
				l.Logf("responding with body from %s", response.Fixture)
				ts.EnqueuedResponses[r.URL.Path] = res[1:]
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
	if _, ok := ts.EnqueuedResponses[path]; !ok {
		ts.EnqueuedResponses[path] = []RecordedResponse{}
	}
	ts.EnqueuedResponses[path] = append(ts.EnqueuedResponses[path], RecordedResponse{
		Body:       body,
		StatusCode: statusCode,
		Fixture:    fixture,
	})
}

func (ts *TestServer) LastRequest() (*http.Request, error) {
	if len(ts.ObservedRequests) < 1 {
		return nil, errors.New("no request")
	}
	return ts.ObservedRequests[len(ts.ObservedRequests)-1], nil
}

func NewTestServer(l logging.SimpleLogger) (*TestServer, func()) {
	ts := &TestServer{
		EnqueuedResponses: make(map[string][]RecordedResponse),
	}
	ts.Server = httptest.NewServer(http.HandlerFunc(ts.stub(l)))
	return ts, func() { ts.Server.Close() }
}
