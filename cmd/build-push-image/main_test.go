package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNexusBuildArgs(t *testing.T) {
	tests := map[string]struct {
		nexusUrl           string
		nexusUsername      string
		nexusPassword      string
		baNexusUsername    string
		baNexusPassword    string
		baNexusHost        string
		baNexusAuth        string
		baNexusUrlWithAuth string
	}{
		"simple-password": {
			nexusUrl:           "https://nexus-ods.example.openshiftapps.com",
			nexusUsername:      "un",
			nexusPassword:      "pw",
			baNexusUsername:    "un",
			baNexusPassword:    "pw",
			baNexusHost:        "nexus-ods.example.openshiftapps.com",
			baNexusAuth:        "un:pw",
			baNexusUrlWithAuth: "https://un:pw@nexus-ods.example.openshiftapps.com",
		},
		"simple-username-only": {
			nexusUrl:           "https://nexus-ods.example.openshiftapps.com",
			nexusUsername:      "un",
			nexusPassword:      "",
			baNexusUsername:    "un",
			baNexusPassword:    "",
			baNexusHost:        "nexus-ods.example.openshiftapps.com",
			baNexusAuth:        "un",
			baNexusUrlWithAuth: "https://un@nexus-ods.example.openshiftapps.com",
		},
		"simple-no-auth": {
			nexusUrl:           "https://nexus-ods.example.openshiftapps.com",
			nexusUsername:      "",
			nexusPassword:      "",
			baNexusUsername:    "",
			baNexusPassword:    "",
			baNexusHost:        "nexus-ods.example.openshiftapps.com",
			baNexusAuth:        "",
			baNexusUrlWithAuth: "https://nexus-ods.example.openshiftapps.com",
		},
		"complex-password": {
			nexusUrl:           "https://nexus-ods.example.openshiftapps.com",
			nexusUsername:      "user: mypw-to-follow",
			nexusPassword:      "a secret",
			baNexusUsername:    "user%3A%20mypw-to-follow",
			baNexusPassword:    "a%20secret",
			baNexusHost:        "nexus-ods.example.openshiftapps.com",
			baNexusAuth:        "user%3A%20mypw-to-follow:a%20secret",
			baNexusUrlWithAuth: "https://user%3A%20mypw-to-follow:a%20secret@nexus-ods.example.openshiftapps.com",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			opts := options{
				nexusURL:      tc.nexusUrl,
				nexusUsername: tc.nexusUsername,
				nexusPassword: tc.nexusPassword,
			}
			args, err := nexusBuildArgs(opts)
			if err != nil {
				t.Fatal(err)
			}

			expected := []string{
				fmt.Sprintf("--build-arg=nexusUrl=\"%s\"", tc.nexusUrl),
				fmt.Sprintf("--build-arg=nexusUsername=\"%s\"", tc.baNexusUsername),
				fmt.Sprintf("--build-arg=nexusPassword=\"%s\"", tc.baNexusPassword),
				fmt.Sprintf("--build-arg=nexusHost=\"%s\"", tc.baNexusHost),
				fmt.Sprintf("--build-arg=nexusAuth=\"%s\"", tc.baNexusAuth),
				fmt.Sprintf("--build-arg=nexusUrlWithAuth=\"%s\"", tc.baNexusUrlWithAuth),
			}
			if diff := cmp.Diff(expected, args); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}
}
