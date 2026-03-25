package api

import (
	"net/http"
	"os"
	"strings"
)

// TokenProvider devolve o token GitHub que deve ser usado para uma request.
type TokenProvider interface {
	TokenForRequest(r *http.Request) string
}

// simpleTokenProvider implementa a lógica:
// 1) tenta header X-GitHub-Token
// 2) fallback para GITHUB_TOKEN do ambiente
type simpleTokenProvider struct {
	defaultToken string
}

func NewSimpleTokenProvider() TokenProvider {
	return &simpleTokenProvider{
		defaultToken: os.Getenv("GITHUB_TOKEN"),
	}
}

func (p *simpleTokenProvider) TokenForRequest(r *http.Request) string {
	// 1) PAT do usuário via header
	if hdr := r.Header.Get("X-GitHub-Token"); hdr != "" {
		return strings.TrimSpace(hdr)
	}

	// 2) fallback: token global (pode ser vazio; aí chamadas ao GitHub falham)
	return p.defaultToken
}
