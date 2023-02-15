package artifact

import (
	"encoding/json"
	"fmt"
	"os"
)

type Image struct {
	// Ref is the full imageRef including Registry/Repository/Name
	Ref string `json:"image"`
	// Registry is  the registry host
	Registry string `json:"registry"`
	// Repository is the namespace in the context of OpenShift
	Repository string `json:"repository"`
	// Name is the image without Registry and Repository
	Name string `json:"name"`
	// Tag is the git commit SHA
	Tag string `json:"tag"`
	// Digest is the SHA of the image
	Digest string `json:"digest"`
}

func (i Image) ImageName() string {
	return fmt.Sprintf("%s:%s", i.Name, i.Tag)
}

func ReadFromFile(filename string) (*Image, error) {
	artifactContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", filename, err)
	}
	var img Image
	err = json.Unmarshal(artifactContent, &img)
	if err != nil {
		return nil, fmt.Errorf(
			"unmarshal %s: %w.\nFile content:\n%s",
			filename, err, string(artifactContent),
		)
	}
	return &img, nil
}
