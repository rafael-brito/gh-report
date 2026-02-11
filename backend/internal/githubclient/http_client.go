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
	// aqui entra o cache em memória depois
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

// Implementações vazias por enquanto; depois preenchemos
func (c *httpClient) ListCommitsByFile(ctx context.Context, params ListCommitsByFileParams) ([]Commit, error) {
	return nil, nil
}

func (c *httpClient) ListPRsByCommit(ctx context.Context, owner, repo, sha string) ([]PRShort, error) {
	return nil, nil
}

func (c *httpClient) CompareCommits(ctx context.Context, params CompareParams) ([]Commit, error) {
	return nil, nil
}

func (c *httpClient) GetPRByNumber(ctx context.Context, owner, repo string, number int) (*PRShort, error) {
	return nil, nil
}
