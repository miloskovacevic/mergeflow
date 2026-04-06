package gitlab

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/miloskovacevic/mergeflow/internal/domain"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	HttpClient *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HttpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type GitlabUser struct {
	Username string `json:"username"`
}
type gitlabMergeRequest struct {
	ID          int        `json:"id"`
	IID         int        `json:"iid"`
	ProjectID   int        `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	MergedAt    *time.Time `json:"merged_at"`
	Author      GitlabUser `json:"author"`
}

func (c *Client) GetMergeRequests(ctx context.Context, filter domain.MergeRequestFilter) (domain.MergeRequestCollection, error) {
	endpoint := fmt.Sprintf("%sprojects/%d/merge_requests", c.BaseURL, filter.ProjectId)

	params := url.Values{}
	if filter.Status != "" {
		params.Set("state", filter.Status)
	}
	//if filter.TargetBranch != "" {
	//	params.Set("target_branch", filter.TargetBranch)
	//}
	//params.Set("per_page", "100")

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var gitlabMRs []gitlabMergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&gitlabMRs); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	collection := make(domain.MergeRequestCollection, 0, len(gitlabMRs))
	for _, gmr := range gitlabMRs {
		collection = append(collection, domain.MergeRequest{
			ID:        strconv.Itoa(gmr.ID),
			Title:     gmr.Title,
			CreatedAt: gmr.CreatedAt,
			Author:    gmr.Author.Username,
			Status:    gmr.State,
		})
	}

	return collection, nil
}
