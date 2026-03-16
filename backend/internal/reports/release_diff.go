package reports

import (
	"context"
	"strings"
	"time"

	"github.com/rafael-brito/gh-report/backend/internal/githubclient"
)

type ReleaseDiffService interface {
	GetReleaseDiffReport(ctx context.Context, params ReleaseDiffParams) (*ReleaseDiffReport, error)
}

type releaseDiffService struct {
	gh githubclient.Client
}

const defaultMaxCompareCommits = 250 // limite de segurança; pode ajustar

func NewReleaseDiffService(gh githubclient.Client) ReleaseDiffService {
	return &releaseDiffService{gh: gh}
}

func (s *releaseDiffService) GetReleaseDiffReport(ctx context.Context, params ReleaseDiffParams) (*ReleaseDiffReport, error) {
	now := time.Now().UTC()

	// 1. Comparar commits entre From e To
	commits, err := s.gh.CompareCommits(ctx, githubclient.CompareParams{
		RepoOwner:  params.Repo.Owner,
		RepoName:   params.Repo.Name,
		Base:       params.From,
		Head:       params.To,
		MaxCommits: defaultMaxCompareCommits,
	})
	if err != nil {
		return nil, err
	}

	// Se não houver commits, retornamos um relatório vazio, mas ainda assim informativo
	if len(commits) == 0 {
		return &ReleaseDiffReport{
			Repository:  params.Repo,
			FromTag:     params.From,
			ToTag:       params.To,
			GeneratedAt: now,
			Summary: ReleaseSummary{
				TotalPRs: 0,
				ByType:   make(map[ReleasePRType]int),
			},
			PRs: []ReleasePR{},
		}, nil
	}

	// 2. Para cada commit, buscar PRs associadas
	prMap := map[int]*ReleasePR{}

	for _, cmt := range commits {
		prs, err := s.gh.ListPRsByCommit(ctx, params.Repo.Owner, params.Repo.Name, cmt.SHA)
		if err != nil {
			return nil, err
		}
		if len(prs) == 0 {
			continue
		}

		for _, ghPR := range prs {
			if _, exists := prMap[ghPR.Number]; exists {
				continue
			}

			labels := ghPR.Labels

			rp := &ReleasePR{
				Number:             ghPR.Number,
				Title:              ghPR.Title,
				URL:                ghPR.URL,
				AuthorLogin:        ghPR.AuthorLogin,
				AuthorAvatarURL:    ghPR.AuthorAvatarURL,
				MergedAt:           ghPR.MergedAt,
				Labels:             labels,
				TypeClassification: ReleasePRTypeUnknown, // será classificado abaixo
			}

			rp.TypeClassification = classifyReleasePRType(labels)
			prMap[ghPR.Number] = rp
		}
	}

	// 3. Converter map em slice e montar summary
	prs := make([]ReleasePR, 0, len(prMap))
	summary := ReleaseSummary{
		TotalPRs: len(prMap),
		ByType:   map[ReleasePRType]int{},
	}

	for _, pr := range prMap {
		prs = append(prs, *pr)
		summary.ByType[pr.TypeClassification]++
	}

	report := &ReleaseDiffReport{
		Repository:  params.Repo,
		FromTag:     params.From,
		ToTag:       params.To,
		GeneratedAt: now,
		Summary:     summary,
		PRs:         prs,
	}

	return report, nil
}

// classifyReleasePRType aplica uma heurística simples baseada em labels
func classifyReleasePRType(labels []string) ReleasePRType {
	if len(labels) == 0 {
		return ReleasePRTypeUnknown
	}

	lower := make([]string, len(labels))
	for i, l := range labels {
		lower[i] = strings.ToLower(l)
	}

	has := func(substr string) bool {
		for _, l := range lower {
			if strings.Contains(l, substr) {
				return true
			}
		}
		return false
	}

	switch {
	case has("bug"), has("bugfix"), has("hotfix"):
		return ReleasePRTypeBugfix
	case has("feature"), has("enhancement"):
		return ReleasePRTypeFeature
	case has("chore"), has("refactor"), has("technical-debt"), has("techdebt"):
		return ReleasePRTypeTechnical
	case has("improvement"):
		return ReleasePRTypeImprovement
	default:
		return ReleasePRTypeUnknown
	}
}
