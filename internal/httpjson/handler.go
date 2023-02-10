package httpjson

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler is an HTTP handler implementing http.Handler.
type Handler func(w http.ResponseWriter, r *http.Request) (any, error)

// ServeHTTP implements http.Handler.
// If an error is returned from h, it is converted to a JSON error.
// Otherwise, the returned value is JSON encoded.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err := h(w, r)
	if err != nil {
		var pd ProblemDetailer
		if pe, ok := err.(ProblemDetailer); ok {
			pd = pe
		} else {
			pd = &StatusProblem{Status: http.StatusInternalServerError, Err: err}
		}
		log.Println(pd)
		JSONError(w, pd.ProblemDetail(), pd.Code())
		return
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err)
		JSONError(w, "internal server error", http.StatusInternalServerError)
	}
}

// JSONError is https://pkg.go.dev/net/http#Error, but for JSON.
func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if e := json.NewEncoder(w).Encode(err); e != nil {
		log.Println(e)
	}
}
