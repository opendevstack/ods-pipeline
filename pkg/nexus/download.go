package nexus

import (
	"io"
	"net/http"
	"os"
	"path"
)

func (c *Client) Download(url, outfile string) (int64, error) {
	if len(outfile) == 0 {
		outfile = path.Base(url)
	}
	out, err := os.Create(outfile)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	// TODO: timeout
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(c.Username, c.Password))
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}
