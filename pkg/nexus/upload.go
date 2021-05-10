package nexus

import (
	"fmt"
	"os"
	"path/filepath"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

// Upload uploads
func Upload(URL, user, password, repository, group, file string) error {
	rm, err := nexusrm.New(
		URL,
		user,
		password,
	)
	if err != nil {
		return fmt.Errorf("could not create nexus client: %w", err)
	}
	// group has leading slash
	link := fmt.Sprintf("%s/repository/%s%s/%s", URL, repository, group, file)
	fmt.Println("Uploading", file, "to", link)

	osFile, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", file, err)
	}

	filename := filepath.Base(file)

	uploadAssetRaw := nexusrm.UploadAssetRaw{
		File:     osFile,
		Filename: filename,
	}
	uploadComponentRaw := nexusrm.UploadComponentRaw{
		Directory: group,
		Tag:       "",
		Assets:    []nexusrm.UploadAssetRaw{uploadAssetRaw},
	}
	err = nexusrm.UploadComponent(rm, repository, uploadComponentRaw)
	if err != nil {
		return fmt.Errorf("could not upload component: %w", err)
	}
	return nil
}
