package nexus

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"

	nexusrm "github.com/sonatype-nexus-community/gonexus/rm"
)

type Client struct {
	RM         nexusrm.RM
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

	return &Client{RM: rm, Username: user, Password: password, Repository: repository}, nil
}

// URLs gets URLs
func (c *Client) URLs(group string) ([]string, error) {
	query := nexusrm.NewSearchQueryBuilder().Repository(c.Repository).Group(group)
	assets, err := nexusrm.SearchAssets(c.RM, query)
	if err != nil {
		return nil, fmt.Errorf("could not search assets: %w", err)
	}

	res := []string{}
	for _, a := range assets {
		res = append(res, a.DownloadURL)
	}
	return res, nil
}

// Upload uploads
func (c *Client) Upload(group, file string) error {

	link := fmt.Sprintf("%s/repository/%s%s/%s", c.RM.Info().Host, c.Repository, group, file)
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
	err = nexusrm.UploadComponent(c.RM, c.Repository, uploadComponentRaw)
	if err != nil {
		return fmt.Errorf("could not upload component: %w", err)
	}
	return nil
}

func (c *Client) Download(url string) (int64, error) {
	outfile := path.Base(url)
	out, err := os.Create(outfile)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// TODO: timeout
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
