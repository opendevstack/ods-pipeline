package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func TestGetODSConfig(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *ODS
		wantErr bool
	}{
		{
			name: "read config from ods.yml",
			args: args{
				filename: "testdata/ods.yml",
			},
			want: &ODS{
				Phases: Phases{
					Build: []v1beta1.PipelineTask{
						{
							Name: "backend-build-go",
							Params: []v1beta1.Param{
								{
									Name: "go-image",
									Value: v1beta1.ArrayOrString{
										StringVal: "localhost:5000/ods/ods-go-toolset:latest",
										Type:      "string",
									},
								},
								{
									Name: "sonar-image",
									Value: v1beta1.ArrayOrString{
										StringVal: "localhost:5000/ods/ods-sonar:latest",
										Type:      "string",
									},
								},
								{
									Name: "go-os",
									Value: v1beta1.ArrayOrString{
										StringVal: "linux",
										Type:      "string",
									},
								},
								{
									Name: "go-arch",
									Value: v1beta1.ArrayOrString{
										StringVal: "amd64",
										Type:      "string",
									},
								},
							},
							TaskRef: &v1beta1.TaskRef{
								Kind: v1beta1.ClusterTaskKind,
								Name: "ods-build-go-v0-1-0",
							},
							Workspaces: []v1beta1.WorkspacePipelineTaskBinding{
								{
									Name:      "source",
									Workspace: "shared-workspace",
								}},
						},
						{
							Name: "backend-build-image",
							TaskRef: &v1beta1.TaskRef{
								Kind: v1beta1.ClusterTaskKind,
								Name: "ods-build-image-v0-1-0",
							},
							RunAfter: []string{"backend-build-go"},
							Params: []v1beta1.Param{
								{
									Name: "registry",
									Value: v1beta1.ArrayOrString{
										StringVal: "kind-registry.kind:5000",
										Type:      "string",
									},
								},
								{
									Name: "builder-image",
									Value: v1beta1.ArrayOrString{
										StringVal: "localhost:5000/ods/ods-buildah:latest",
										Type:      "string",
									},
								},
								{
									Name: "tls-verify",
									Value: v1beta1.ArrayOrString{
										StringVal: "false",
										Type:      "string",
									},
								},
							},
							Workspaces: []v1beta1.WorkspacePipelineTaskBinding{
								{
									Name:      "source",
									Workspace: "shared-workspace",
								},
							},
						},
					},
					Deploy: []v1beta1.PipelineTask{
						{
							Name: "backend-deploy",
							TaskRef: &v1beta1.TaskRef{
								Kind: v1beta1.ClusterTaskKind,
								Name: "ods-deploy-helm-v0-1-0",
							},
							Params: []v1beta1.Param{
								{
									Name: "release-name",
									Value: v1beta1.ArrayOrString{
										StringVal: "backend",
										Type:      "string",
									},
								},
								{
									Name: "image",
									Value: v1beta1.ArrayOrString{
										StringVal: "localhost:5000/ods/ods-helm:latest",
										Type:      "string",
									},
								},
							},
							Workspaces: []v1beta1.WorkspacePipelineTaskBinding{
								{
									Name:      "source",
									Workspace: "shared-workspace",
								},
							},
						},
					},
				},
				Repositories: []Repository{
					{Name: "foo", Branch: "master"},
					{Name: "bar", Branch: "master"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetODSConfig(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetODSConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Diff: %s\n", diff)
			}
		})
	}
}
