package manager

import (
	"net/http"
	"testing"
)

const testWebhookSecret = "s3cr3t"

func TestValidatePayload(t *testing.T) {
	body := []byte("hello world")
	correctBodyHMAC := hmacHeader(t, testWebhookSecret, body)
	header := func(name, value string) http.Header {
		h := http.Header{}
		h.Set(name, value)
		return h
	}
	tests := map[string]struct {
		header      http.Header
		payload     string
		secretToken []byte
		wantErr     bool
	}{
		"missing header": {
			header:      http.Header{},
			payload:     "pvc",
			secretToken: []byte(testWebhookSecret),
			wantErr:     true,
		},
		"proper header": {
			header:      header(signatureHeader, correctBodyHMAC),
			payload:     string(body),
			secretToken: []byte(testWebhookSecret),
			wantErr:     false,
		},
		"lowercase header": {
			header:      header("x-hub-signature", correctBodyHMAC),
			payload:     string(body),
			secretToken: []byte(testWebhookSecret),
			wantErr:     false,
		},
		"empty secret token": {
			header:      header(signatureHeader, correctBodyHMAC),
			payload:     string(body),
			secretToken: []byte(""),
			wantErr:     true,
		},
		"wrong secret token": {
			header:      header(signatureHeader, correctBodyHMAC),
			payload:     string(body),
			secretToken: []byte("abc"),
			wantErr:     true,
		},
		"incorrect signature": {
			header:      header(signatureHeader, "wrong"),
			payload:     string(body),
			secretToken: []byte(testWebhookSecret),
			wantErr:     true,
		},
		"missing signature": {
			header:      header(signatureHeader, ""),
			payload:     string(body),
			secretToken: []byte(testWebhookSecret),
			wantErr:     true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := validatePayload(tc.header, body, tc.secretToken)
			if tc.wantErr != (err != nil) {
				t.Fatalf("want err: %v, got: %v", tc.wantErr, err)
			}
		})
	}
}
