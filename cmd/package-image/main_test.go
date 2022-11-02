package main

import (
	"testing"
)

func TestAquaScanURL(t *testing.T) {
	tests := map[string]struct {
		aquaURL string
	}{
		"base URL without trailing slash": {
			aquaURL: "https://console.example.com",
		},
		"base URL with trailing slash": {
			aquaURL: "https://console.example.com/",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := options{aquaURL: tc.aquaURL, aquaRegistry: "ods"}
			u, err := aquaScanURL(opts, "foo")
			if err != nil {
				t.Fatal(err)
			}
			want := "https://console.example.com/#/images/ods/foo/vulns"
			if u != want {
				t.Fatalf("want: %s, got: %s", want, u)
			}
		})
	}
}
