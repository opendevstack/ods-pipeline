package manager

import (
	"errors"
	"fmt"
	"net/http"

	// The github package is used here because the Bitbucket interceptor of
	// Tekton Triggers uses it to validate incoming webhook requests, see
	// https://github.com/tektoncd/triggers/tree/main/pkg/interceptors/bitbucket.
	github "github.com/google/go-github/v42/github"
)

const signatureHeader = "X-Hub-Signature"

// Canonical updates the map keys to use the Canonical name
func canonicalHeader(h map[string][]string) http.Header {
	c := map[string][]string{}
	for k, v := range h {
		c[http.CanonicalHeaderKey(k)] = v
	}
	return http.Header(c)
}

// validatePayload errors if the payload does not match the signature provided
// in the header. The secretToken is shared with Bitbucket.
func validatePayload(h http.Header, payload, secretToken []byte) error {
	headers := canonicalHeader(h)
	signature := headers.Get(signatureHeader)
	if signature == "" {
		return fmt.Errorf("no %s set", signatureHeader)
	}
	if len(secretToken) == 0 {
		return errors.New("refuse to validate with empty secret")
	}
	if err := github.ValidateSignature(signature, payload, secretToken); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
