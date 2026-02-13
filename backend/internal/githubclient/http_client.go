package githubclient

import (
	"context"
	"net/http"
	"time"
)

type httpClient struct {
	http    *http.Client
	token   string
	baseURL string
	// cache entra depois
}

func NewHTTPClient(token string) Client {
	return &httpClient{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		token:   token,
		baseURL: "https://api.github.com",
	}
}

// Stubs (alguns ainda vazios, vamos preencher depois)

func (c *httpClient) ListPRsByCommit(ctx context.Context, owner, repo, sha string) ([]PRShort, error) {
	return nil, nil
}

func (c *httpClient) CompareCommits(ctx context.Context, params CompareParams) ([]Commit, error) {
	return nil, nil
}

func (c *httpClient) GetPRByNumber(ctx context.Context, owner, repo string, number int) (*PRShort, error) {
	return nil, nil
}
