package githubclient

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type httpClient struct {
	http    *http.Client
	token   string
	baseURL string

	mu               sync.Mutex
	prsByCommitCache map[string][]PRShort // chave: owner/repo/sha
	prDetailsCache   map[string]*PRShort  // chave: owner/repo/number
}

func NewHTTPClient(token string) Client {
	return &httpClient{
		http: &http.Client{
			Timeout: 15 * time.Second,
		},
		token:            token,
		baseURL:          "https://api.github.com",
		prsByCommitCache: make(map[string][]PRShort),
		prDetailsCache:   make(map[string]*PRShort),
	}
}

func (c *httpClient) GetPRByNumber(ctx context.Context, owner, repo string, number int) (*PRShort, error) {
	return nil, nil
}
