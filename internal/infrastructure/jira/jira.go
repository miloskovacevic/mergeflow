package jira

type Client struct {
	BaseUrl string
	Token   string
}

func NewClient(baseUrl, token string) *Client {
	return &Client{
		BaseUrl: baseUrl,
		Token:   token,
	}
}

func (c *Client) GetProject(project string) string {
	return "fetching project  " + project + " from " + c.BaseUrl
}
