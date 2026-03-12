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

	switch params.Mode {
	case FileHistoryModePRs:
		return s.buildPRModeReport(ctx, params, commits)
	default:
		// fallback para commits
		return s.buildCommitsModeReport(params, commits), nil
	}
}

// ---- modo commits (já funcionava, só extraí para função separada) ----

func (s *fileHistoryService) buildCommitsModeReport(params FileHistoryParams, commits []githubclient.Commit) *FileHistoryReport {
	entries := make([]FileHistoryEntry, 0, len(commits))
	authorCount := map[string]int{}

	for _, c := range commits {
		if c.AuthorLogin != "" {
			authorCount[c.AuthorLogin]++
		}

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
		topAuthors = append(topAuthors, TopAuthorStat{
			Login:   login,
			Commits: count,
		})
	}

	stats := FileHistoryStats{
		TotalCommits: len(commits),
		TotalPRs:     0,
		TopAuthors:   topAuthors,
	}

	return &FileHistoryReport{
		Repository:  params.Repo,
		FilePath:    params.File,
		Mode:        FileHistoryModeCommits,
		Limit:       params.Limit,
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
		Stats:       stats,
	}
}

// ---- modo PRs ----

func (s *fileHistoryService) buildPRModeReport(ctx context.Context, params FileHistoryParams, commits []githubclient.Commit) (*FileHistoryReport, error) {
	// map: prNumber -> Release-like PR info + commits
	type prAgg struct {
		pr      githubclient.PRShort
		commits []FileHistoryCommit
	}

	prMap := map[int]*prAgg{}
	authorCount := map[string]int{}
	totalCommits := 0

	for _, cmt := range commits {
		totalCommits++
		if cmt.AuthorLogin != "" {
			authorCount[cmt.AuthorLogin]++
		}

		// Buscar PRs associadas a este commit
		prs, err := s.gh.ListPRsByCommit(ctx, params.Repo.Owner, params.Repo.Name, cmt.SHA)
		if err != nil {
			return nil, err
		}

		// Se não há PRs, podemos ou:
		// - Ignorar esses commits (histórico só por PR)
		// - Incluir como "commit solto" (tipo extra)
		// Para agora, vamos ignorar commit sem PR no modo PRs.
		if len(prs) == 0 {
			continue
		}

		// Normalmente haverá 1 PR por commit nesse endpoint
		for _, pr := range prs {
			agg, ok := prMap[pr.Number]
			if !ok {
				agg = &prAgg{
					pr:      pr,
					commits: []FileHistoryCommit{},
				}
				prMap[pr.Number] = agg
			}

			agg.commits = append(agg.commits, FileHistoryCommit{
				SHA:          cmt.SHA,
				Message:      cmt.Message,
				URL:          cmt.URL,
				AuthorLogin:  cmt.AuthorLogin,
				AuthorAvatar: cmt.AuthorAvatarURL,
				CommittedAt:  cmt.CommittedAt,
			})
		}
	}

	entries := make([]FileHistoryEntry, 0, len(prMap))
	for _, agg := range prMap {
		// Encontrar data de ordenação (ex: mergedAt ou data do último commit)
		var orderTs time.Time
		if agg.pr.MergedAt != nil {
			orderTs = *agg.pr.MergedAt
		} else if len(agg.commits) > 0 {
			orderTs = agg.commits[0].CommittedAt
		} else {
			orderTs = time.Now().UTC()
		}

		entries = append(entries, FileHistoryEntry{
			Type: FileHistoryEntryTypePR,
			PR: &FileHistoryPREntry{
				PRNumber:    agg.pr.Number,
				PRTitle:     agg.pr.Title,
				PRURL:       agg.pr.URL,
				PRMergedAt:  agg.pr.MergedAt,
				Commits:     agg.commits,
				AuthorLogin: agg.pr.AuthorLogin,
			},
			OrderTs: orderTs,
		})
	}

	topAuthors := make([]TopAuthorStat, 0, len(authorCount))
	for login, count := range authorCount {
		topAuthors = append(topAuthors, TopAuthorStat{
			Login:   login,
			Commits: count,
		})
	}

	stats := FileHistoryStats{
		TotalCommits: totalCommits,
		TotalPRs:     len(prMap),
		TopAuthors:   topAuthors,
	}

	report := &FileHistoryReport{
		Repository:  params.Repo,
		FilePath:    params.File,
		Mode:        FileHistoryModePRs,
		Limit:       params.Limit,
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
		Stats:       stats,
	}

	return report, nil
}
