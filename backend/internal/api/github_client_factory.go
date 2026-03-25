package api

import (
	"sync"

	"github.com/rafael-brito/gh-report/backend/internal/githubclient"
)

// GitHubClientFactory cria/reusa clients com base no token fornecido.
type GitHubClientFactory interface {
	ForToken(token string) githubclient.Client
}

type cachedGitHubClientFactory struct {
	mu      sync.Mutex
	clients map[string]githubclient.Client
}

func NewGitHubClientFactory() GitHubClientFactory {
	return &cachedGitHubClientFactory{
		clients: make(map[string]githubclient.Client),
	}
}

func (f *cachedGitHubClientFactory) ForToken(token string) githubclient.Client {
	f.mu.Lock()
	defer f.mu.Unlock()

	if c, ok := f.clients[token]; ok {
		return c
	}

	c := githubclient.NewHTTPClient(token)
	f.clients[token] = c
	return c
}
