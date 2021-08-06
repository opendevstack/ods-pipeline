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

	c.logger().Debugf("Download %s to %s", url, outfile)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Basic "+c.basicAuth())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}
