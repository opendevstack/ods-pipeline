package nexus

import (
	"fmt"
	"os"
	"path/filepath"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Upload a file to a repository group
func (c *Client) Upload(group, file string) error {

	filename := filepath.Base(file)

	link := fmt.Sprintf("%s/repository/%s%s/%s", c.rm.Info().Host, c.Repository(), group, filename)
	c.logger().Debugf("Uploading %s to %s", file, link)

	osFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", file, err)
	}

	uploadAssetRaw := nexusrm.UploadAssetRaw{
		File:     osFile,
		Filename: filename,
	}
	uploadComponentRaw := nexusrm.UploadComponentRaw{
		Directory: group,
		Tag:       "",
		Assets:    []nexusrm.UploadAssetRaw{uploadAssetRaw},
	}
	err = nexusrm.UploadComponent(c.rm, c.Repository(), uploadComponentRaw)
	if err != nil {
		return fmt.Errorf("could not upload component: %w", err)
	}
	return nil
}
