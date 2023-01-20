package artifact

import "fmt"

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
