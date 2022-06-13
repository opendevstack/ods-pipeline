package tekton

import "testing"

func TestPipelineRunURL(t *testing.T) {
	tests := map[string]struct {
		consoleURL string
	}{
		"base URL without trailing slash": {
			consoleURL: "https://console.example.com",
		},
		"base URL with trailing slash": {
			consoleURL: "https://console.example.com/",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := PipelineRunURL(tc.consoleURL, "foo", "bar-ab12c")
			if err != nil {
				t.Fatal(err)
			}
			want := "https://console.example.com/k8s/ns/foo/tekton.dev~v1beta1~PipelineRun/bar-ab12c/"
			if u != want {
				t.Fatalf("want: %s, got: %s", want, u)
			}
		})
	}
}
