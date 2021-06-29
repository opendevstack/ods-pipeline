package config

import (
	"reflect"
	"testing"

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
					Build: Phase{
						RunPolicy: "Serial",
						Tasks: []v1beta1.PipelineTask{
							{
								Name: "ods-build-go-v0-1-0",
								Params: []v1beta1.Param{
									{
										Name: "go-image",
										Value: v1beta1.ArrayOrString{
											StringVal: "localhost:5000/ods/ods-go-toolset:latest",
										},
									},
								},
							},
						},
					},
					Deploy: Phase{
						RunPolicy: "Serial",
						Tasks: []v1beta1.PipelineTask{
							{
								Name: "ods-build-image-v0-1-0",
							},
						},
					},
				},
				Repositories: []Repository{
					{Name: "foo", URL: "http://foo"},
					{Name: "bar", URL: "http://bar"},
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetODSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
