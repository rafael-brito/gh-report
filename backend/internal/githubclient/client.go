package githubclient

import (
	"context"
	"time"
)

type Commit struct {
	SHA             string
	Message         string
	AuthorLogin     string
	AuthorAvatarURL string
	CommittedAt     time.Time
	URL             string
}

type PRShort struct {
	Number          int
	Title           string
	URL             string
	AuthorLogin     string
	AuthorAvatarURL string
	MergedAt        *time.Time
	Labels          []string
}

// Params

type ListCommitsByFileParams struct {
	RepoOwner string
	RepoName  string
	FilePath  string
	Limit     int
}

type CompareParams struct {
	RepoOwner  string
	RepoName   string
	Base       string
	Head       string
	MaxCommits int // para evitar explosão
}

type Client interface {
	// Autenticação já embutida via token no construtor

	// Relatório 1
	ListCommitsByFile(ctx context.Context, params ListCommitsByFileParams) ([]Commit, error)
	ListPRsByCommit(ctx context.Context, repoOwner, repoName, sha string) ([]PRShort, error)

	// Relatório 2
	CompareCommits(ctx context.Context, params CompareParams) ([]Commit, error)
	GetPRByNumber(ctx context.Context, repoOwner, repoName string, number int) (*PRShort, error)
}
