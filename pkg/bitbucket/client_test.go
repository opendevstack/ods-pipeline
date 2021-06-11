package bitbucket

import (
	"net/http/httptest"
	"testing"

	"github.com/opendevstack/pipeline/internal/serverstub"
)

func testClient(t *testing.T, endpoints map[string]*serverstub.FakeResponse) *Client {
	st, err := serverstub.Stub(t, endpoints)
	if err != nil {
		t.Fatal(err)
	}
	fakeBitbucket := httptest.NewServer(st)
	return NewClient(&ClientConfig{
		APIToken: "s3cr3t", // does not matter
		BaseURL:  fakeBitbucket.URL,
	})
}
