package httpjson

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	ts := httptest.NewServer(Handler(func(w http.ResponseWriter, r *http.Request) (any, error) {
		return nil, NewInternalProblem("foo", nil)
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	want := `{"title":"Internal Server Error","detail":"foo"}`
	if strings.TrimSpace(string(greeting)) != want {
		t.Fatal("not euqal, got", string(greeting))
	}
}
