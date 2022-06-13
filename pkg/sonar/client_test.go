package sonar

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, baseURL string) *Client {
	c, err := NewClient(&ClientConfig{BaseURL: baseURL})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestGetRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	}))
	defer srv.Close()
	tests := map[string]struct {
		baseURL string
	}{
		"base URL without trailing slash": {
			baseURL: srv.URL,
		},
		"base URL with trailing slash": {
			baseURL: srv.URL + "/",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			bitbucketClient := testClient(t, tc.baseURL)
			requestPath := "/foo"
			code, out, err := bitbucketClient.get(requestPath)
			if err != nil {
				t.Fatal(err)
			}
			if code != 200 {
				t.Fatal("expected 200")
			}
			if string(out) != requestPath {
				t.Fatalf("expected %s, got: %s", requestPath, string(out))
			}
		})
	}
}
