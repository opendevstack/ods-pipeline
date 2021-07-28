package artifact

type Image struct {
	Image      string `json:"image"`
	Registry   string `json:"registry"`
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	Digest     string `json:"digest"`
}
