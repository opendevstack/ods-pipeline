package sonar

func testClient(baseURL string) *Client {
	return NewClient(&ClientConfig{BaseURL: baseURL})
}
