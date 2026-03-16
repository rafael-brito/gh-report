package githubclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ghCompareCommit struct {
	SHA     string          `json:"sha"`
	HTMLURL string          `json:"html_url"`
	Commit  ghCommitInfo    `json:"commit"`
	Author  *ghCommitAuthor `json:"author"` // pode ser null
}

type ghCompareResponse struct {
	Commits []ghCompareCommit `json:"commits"`
}

// CompareCommits chama o endpoint de comparação de commits do GitHub:
// GET /repos/{owner}/{repo}/compare/{base}...{head}
func (c *httpClient) CompareCommits(ctx context.Context, params CompareParams) ([]Commit, error) {
	if params.MaxCommits <= 0 {
		// define um default razoável para evitar explosão de dados
		params.MaxCommits = 250
	}

	u, err := url.Parse(fmt.Sprintf(
		"%s/repos/%s/%s/compare/%s...%s",
		c.baseURL,
		params.RepoOwner,
		params.RepoName,
		params.Base,
		params.Head,
	))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("github CompareCommits: base or head not found (status 404)")
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github CompareCommits: status %d", resp.StatusCode)
	}

	var ghResp ghCompareResponse
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, err
	}

	out := make([]Commit, 0, len(ghResp.Commits))
	for i, gc := range ghResp.Commits {
		if i >= params.MaxCommits {
			break
		}

		authorLogin := ""
		authorAvatar := ""
		if gc.Author != nil {
			authorLogin = gc.Author.Login
			authorAvatar = gc.Author.AvatarURL
		}

		out = append(out, Commit{
			SHA:             gc.SHA,
			Message:         gc.Commit.Message,
			AuthorLogin:     authorLogin,
			AuthorAvatarURL: authorAvatar,
			CommittedAt:     gc.Commit.Author.Date,
			URL:             gc.HTMLURL,
		})
	}

	return out, nil
}
