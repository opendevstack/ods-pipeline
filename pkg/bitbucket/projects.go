package bitbucket

type Project struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
	Links       struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}
