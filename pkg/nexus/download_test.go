package nexus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/opendevstack/pipeline/internal/testfile"
	"github.com/opendevstack/pipeline/pkg/logging"
)

func TestDownload(t *testing.T) {
	want := "Hello, client"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth := r.Header.Get("Authorization")
		wantAuth := "Basic dXNlcm5hbWU6cGFzc3dvcmQ=" // username:password
		if gotAuth != wantAuth {
			t.Fatalf("want basic auth: %s, got: %s", wantAuth, gotAuth)
		}
		fmt.Fprint(w, want)
	}))
	defer ts.Close()

	nexusClient, err := NewClient(&ClientConfig{
		BaseURL:  "http://localhost",
		Username: "username",
		Password: "password",
		Logger:   &logging.LeveledLogger{Level: logging.LevelDebug},
	})
	if err != nil {
		t.Fatal(err)
	}

	outfile := "res.out"
	defer func() {
		if _, err := os.Stat(outfile); err == nil {
			os.Remove(outfile)
		}
	}()
	_, err = nexusClient.Download(ts.URL, outfile)
	if err != nil {
		t.Fatal(err)
	}
	got := string(testfile.ReadFileOrFatal(t, outfile))
	if want != got {
		t.Fatalf("want: %s, got: %s", want, got)
	}
}
