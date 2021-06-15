package sonar

import (
	"net/http"
	"net/http/httptest"
)

func testClient(baseURL string) *Client {
	return NewClient(&ClientConfig{BaseURL: baseURL})
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
