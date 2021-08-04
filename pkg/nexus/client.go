package nexus

import (
	"encoding/base64"
	"fmt"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

type Client struct {
	RM         nexusrm.RM
	URL        string
	Username   string
	Password   string
	Repository string
}

// NewClient initializes client
func NewClient(URL, user, password, repository string) (*Client, error) {
	rm, err := nexusrm.New(
		URL,
		user,
		password,
	)
	if err != nil {
		return nil, fmt.Errorf("could not create nexus client: %w", err)
	}

	return &Client{RM: rm, URL: URL, Username: user, Password: password, Repository: repository}, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
