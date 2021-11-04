package repository

import (
	"fmt"

	"github.com/opendevstack/pipeline/pkg/bitbucket"
	"github.com/opendevstack/pipeline/pkg/config"
)

// GetODSConfig returns a *config.ODS for given project/repository at gitFullRef.
// If retrieving fails or not ods.y(a)ml file exists, it errors.
func GetODSConfig(bitbucketClient bitbucket.RawClientInterface, project, repository, gitFullRef string) (*config.ODS, error) {
	var body []byte
	var getErr error
	for _, c := range config.ODSFileCandidates {
		b, err := bitbucketClient.RawGet(project, repository, c, gitFullRef)
		if err == nil {
			body = b
			getErr = nil
			break
		}
		getErr = err
	}
	if getErr != nil {
		return nil, fmt.Errorf("could not download ODS config for repo %s: %w", repository, getErr)
	}

	if body == nil {
		return nil, fmt.Errorf("no ODS config located in repo %s", repository)
	}
	return config.Read(body)
}
