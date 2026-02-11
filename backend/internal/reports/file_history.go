package reports

import (
	"context"
	"time"

	"gh-report/internal/githubclient"
)

type FileHistoryService interface {
	GetFileHistoryReport(ctx context.Context, params FileHistoryParams) (*FileHistoryReport, error)
}

type fileHistoryService struct {
	gh githubclient.Client
}

func NewFileHistoryService(gh githubclient.Client) FileHistoryService {
	return &fileHistoryService{gh: gh}
}

func (s *fileHistoryService) GetFileHistoryReport(ctx context.Context, params FileHistoryParams) (*FileHistoryReport, error) {
	// 1. Buscar commits que alteraram o arquivo
	commits, err := s.gh.ListCommitsByFile(ctx, githubclient.ListCommitsByFileParams{
		RepoOwner: params.Repo.Owner,
		RepoName:  params.Repo.Name,
		FilePath:  params.File,
		Limit:     params.Limit,
	})
	if err != nil {
		return nil, err
	}

	// 2. Para cada commit, buscar PRs associadas (com cache interno no githubclient)
	// 3. Montar FileHistoryEntry[] de acordo com params.Mode
	// 4. Calcular stats (total_commits, total_prs, top_authors)
	// 5. Retornar FileHistoryReport

	report := &FileHistoryReport{
		Repository:  params.Repo,
		FilePath:    params.File,
		Mode:        params.Mode,
		Limit:       params.Limit,
		GeneratedAt: time.Now().UTC(),
		// Entries: ...
		// Stats: ...
	}

	return report, nil
}
