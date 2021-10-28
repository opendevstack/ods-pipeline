package bitbucket

import (
	"strings"
	"testing"

	"github.com/opendevstack/pipeline/test/testserver"
)

func TestRawGet(t *testing.T) {
	at := "refs/heads/master"

	srv, cleanup := testserver.NewTestServer(t)
	defer cleanup()
	bitbucketClient := testClient(srv.Server.URL)

	tests := map[string]struct {
		EnqueuedPath       string
		EnqueuedStatusCode int
		EnqueuedFixture    string
		TestProject        string
		TestRepository     string
		TestFile           string
		WantError          bool
		WantBody           string
	}{
		"example.txt": {
			EnqueuedPath:       "/projects/PRJ/repos/my-repo/raw/example.txt",
			EnqueuedStatusCode: 200,
			EnqueuedFixture:    "bitbucket/example.txt",
			TestProject:        "PRJ",
			TestRepository:     "my-repo",
			TestFile:           "example.txt",
			WantError:          false,
			WantBody:           "hello world",
		},
		"wrong file": {
			EnqueuedPath:       "/projects/PRJ/repos/my-repo/raw/example.txt",
			EnqueuedStatusCode: 200,
			EnqueuedFixture:    "bitbucket/example.txt",
			TestProject:        "PRJ",
			TestRepository:     "my-repo",
			TestFile:           "foo.txt",
			WantError:          true,
			WantBody:           "",
		},
		"wrong auth": {
			EnqueuedPath: "/projects/PRJ/repos/my-repo/raw/blank.txt",
			// Bitbucket actually returns 302 redirect to login when authentication
			// is incorrect but srv.EnqueueResponse does not allow us to set
			// headers, so we fake this case with a 401.
			EnqueuedStatusCode: 401,
			EnqueuedFixture:    "bitbucket/blank.txt",
			TestProject:        "PRJ",
			TestRepository:     "my-repo",
			TestFile:           "blank.txt",
			WantError:          true,
			WantBody:           "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			srv.EnqueueResponse(t, tc.EnqueuedPath, tc.EnqueuedStatusCode, tc.EnqueuedFixture)

			r, err := bitbucketClient.RawGet(
				tc.TestProject, tc.TestRepository, tc.TestFile, at,
			)
			if (err == nil) == tc.WantError {
				t.Fatalf("got err %v, want err: %v", err, tc.WantError)
			}
			if tc.WantBody != "" {
				got := strings.TrimSpace(string(r))
				if got != tc.WantBody {
					t.Fatalf("got %s, want %s", got, tc.WantBody)
				}
			}
		})
	}
}
