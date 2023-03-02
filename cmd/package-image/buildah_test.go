package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildahBuildArgs(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dockerDir := filepath.Join(basePath, "docker")
	tests := map[string]struct {
		opts     options
		tag      string
		wantArgs []string
		wantErr  string
	}{
		"with default options": {
			opts: defaultOptions,
			tag:  "foo",
			wantArgs: []string{
				"--storage-driver=vfs", "bud", "--format=oci",
				"--tls-verify=true", "--cert-dir=/etc/containers/certs.d",
				"--no-cache",
				"--file=./Dockerfile", "--tag=foo", dockerDir,
			},
		},
		"with blank tag": {
			opts:    defaultOptions,
			tag:     "",
			wantErr: "tag must not be empty",
		},
		"with incorrect buildah extra args": {
			opts:    func(o options) options { o.buildahBuildExtraArgs = "\\"; return o }(defaultOptions),
			tag:     "foo",
			wantErr: "parse extra args (\\): EOF found after escape character",
		},
		"with Nexus args": {
			opts: func(o options) options {
				o.nexusURL = "http://nexus.example.com"
				o.nexusUsername = "developer"
				o.nexusPassword = "s3cr3t"
				return o
			}(defaultOptions),
			tag: "foo",
			wantArgs: []string{
				"--storage-driver=vfs", "bud", "--format=oci",
				"--tls-verify=true", "--cert-dir=/etc/containers/certs.d",
				"--no-cache",
				"--file=./Dockerfile", "--tag=foo",
				"--build-arg=nexusUrl=http://nexus.example.com",
				"--build-arg=nexusUsername=developer",
				"--build-arg=nexusPassword=s3cr3t",
				"--build-arg=nexusHost=nexus.example.com",
				"--build-arg=nexusAuth=developer:s3cr3t",
				"--build-arg=nexusUrlWithAuth=http://developer:s3cr3t@nexus.example.com",
				dockerDir,
			},
		},
		"with debug on": {
			opts: func(o options) options { o.debug = true; return o }(defaultOptions),
			tag:  "foo",
			wantArgs: []string{
				"--storage-driver=vfs", "bud", "--format=oci",
				"--tls-verify=true", "--cert-dir=/etc/containers/certs.d",
				"--no-cache",
				"--file=./Dockerfile", "--tag=foo", "--log-level=debug", dockerDir,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := packageImage{opts: tc.opts}
			got, err := p.buildahBuildArgs(tc.tag)
			if err != nil {
				if tc.wantErr != err.Error() {
					t.Fatalf("want err: '%s', got err: %s", tc.wantErr, err)
				}
			}
			if diff := cmp.Diff(tc.wantArgs, got); diff != "" {
				t.Fatalf("args mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

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
			p := packageImage{opts: opts}
			args, err := p.nexusBuildArgs()
			if err != nil {
				t.Fatal(err)
			}

			expected := []string{
				fmt.Sprintf("--build-arg=nexusUrl=%s", tc.nexusUrl),
				fmt.Sprintf("--build-arg=nexusUsername=%s", tc.baNexusUsername),
				fmt.Sprintf("--build-arg=nexusPassword=%s", tc.baNexusPassword),
				fmt.Sprintf("--build-arg=nexusHost=%s", tc.baNexusHost),
				fmt.Sprintf("--build-arg=nexusAuth=%s", tc.baNexusAuth),
				fmt.Sprintf("--build-arg=nexusUrlWithAuth=%s", tc.baNexusUrlWithAuth),
			}
			if diff := cmp.Diff(expected, args); diff != "" {
				t.Fatalf("expected (-want +got):\n%s", diff)
			}
		})
	}
}
