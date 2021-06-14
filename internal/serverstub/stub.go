package serverstub

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/opendevstack/pipeline/internal/projectpath"
	"github.com/opendevstack/pipeline/pkg/logging"
)

type FakeResponse struct {
	StatusCode int
	Fixture    string
	Body       []byte
}

func Stub(l logging.SimpleLogger, endpoints map[string]*FakeResponse) (http.HandlerFunc, error) {
	for _, eRes := range endpoints {
		if len(eRes.Fixture) > 0 {
			filename := filepath.Join(projectpath.Root, "test/testdata/fixtures", eRes.Fixture)
			res, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, err
			}
			eRes.Body = res
		} else {
			eRes.Body = []byte("")
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if res, ok := endpoints[r.URL.Path]; ok {
			l.Logf("responding with body from %s", res.Fixture)
			w.WriteHeader(res.StatusCode)
			_, err := w.Write(res.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		l.Logf("no stub response registered for path %s", r.URL.Path)
		http.NotFound(w, r)
	}, nil
}
