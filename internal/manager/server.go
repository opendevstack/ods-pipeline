package manager

import (
	"net/http"

	"github.com/opendevstack/ods-pipeline/internal/httpjson"
)

func HealthHandler() http.Handler {
	return http.HandlerFunc(healthEndpoint)
}

func BitbucketHandler(r *BitbucketWebhookReceiver) http.Handler {
	return httpjson.Handler(r.Handle)
}

func healthEndpoint(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`{"health":"ok"}`))
	if err != nil {
		http.Error(w, `{"health":"error"}`, http.StatusInternalServerError)
		return
	}
}
