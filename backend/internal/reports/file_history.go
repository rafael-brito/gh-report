package reports

import (
	"context"
	"time"

	"github.com/rafael-brito/gh-report/backend/internal/githubclient"
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
	if params.Limit <= 0 {
		params.Limit = 10
	}

	commits, err := s.gh.ListCommitsByFile(ctx, githubclient.ListCommitsByFileParams{
		RepoOwner: params.Repo.Owner,
		RepoName:  params.Repo.Name,
		FilePath:  params.File,
		Limit:     params.Limit,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]FileHistoryEntry, 0, len(commits))
	authorCount := map[string]int{}

	for _, c := range commits {
		authorCount[c.AuthorLogin]++

		entry := FileHistoryEntry{
			Type: FileHistoryEntryTypeCommit,
			Commit: &FileHistoryCommit{
				SHA:          c.SHA,
				Message:      c.Message,
				URL:          c.URL,
				AuthorLogin:  c.AuthorLogin,
				AuthorAvatar: c.AuthorAvatarURL,
				CommittedAt:  c.CommittedAt,
			},
			OrderTs: c.CommittedAt,
		}
		entries = append(entries, entry)
	}

	topAuthors := make([]TopAuthorStat, 0, len(authorCount))
	for login, count := range authorCount {
		if login == "" {
			continue
		}
		topAuthors = append(topAuthors, TopAuthorStat{
			Login:   login,
			Commits: count,
		})
	}

	stats := FileHistoryStats{
		TotalCommits: len(commits),
		TotalPRs:     0, // ainda não calculamos PRs
		TopAuthors:   topAuthors,
	}

	report := &FileHistoryReport{
		Repository:  params.Repo,
		FilePath:    params.File,
		Mode:        params.Mode,
		Limit:       params.Limit,
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
		Stats:       stats,
	}

	return report, nil
}
