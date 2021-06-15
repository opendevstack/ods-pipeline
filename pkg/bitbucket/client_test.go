package bitbucket

func testClient(serverURL string) *Client {
	return NewClient(&ClientConfig{
		APIToken: "s3cr3t", // does not matter
		BaseURL:  serverURL,
	})
}
