package reports

import (
	"context"
	"time"

	"gh-report/internal/githubclient"
)

type ReleaseDiffService interface {
	GetReleaseDiffReport(ctx context.Context, params ReleaseDiffParams) (*ReleaseDiffReport, error)
}

type releaseDiffService struct {
	gh githubclient.Client
}

func NewReleaseDiffService(gh githubclient.Client) ReleaseDiffService {
	return &releaseDiffService{gh: gh}
}

func (s *releaseDiffService) GetReleaseDiffReport(ctx context.Context, params ReleaseDiffParams) (*ReleaseDiffReport, error) {
	// 1. Chamar compare
	// 2. A partir dos commits, descobrir PRs (por /commits/{sha}/pulls e/ou regex em commit message)
	// 3. Buscar detalhes das PRs
	// 4. Classificar por tipo
	// 5. Montar summary e retornar

	report := &ReleaseDiffReport{
		Repository:  params.Repo,
		FromTag:     params.From,
		ToTag:       params.To,
		GeneratedAt: time.Now().UTC(),
		// Summary: ...
		// PRs: ...
	}
	return report, nil
}
