package githubclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ghCommitAuthor struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type ghCommitInfo struct {
	Message string `json:"message"`
	Author  struct {
		Date time.Time `json:"date"`
	} `json:"author"`
}

type ghCommit struct {
	SHA     string          `json:"sha"`
	HTMLURL string          `json:"html_url"`
	Commit  ghCommitInfo    `json:"commit"`
	Author  *ghCommitAuthor `json:"author"` // pode ser null
}

func (c *httpClient) ListCommitsByFile(ctx context.Context, params ListCommitsByFileParams) ([]Commit, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	u, err := url.Parse(fmt.Sprintf("%s/repos/%s/%s/commits", c.baseURL, params.RepoOwner, params.RepoName))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("path", params.FilePath)
	q.Set("per_page", fmt.Sprintf("%d", params.Limit))
	u.RawQuery = q.Encode()

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

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github ListCommitsByFile: status %d", resp.StatusCode)
	}

	var ghCommits []ghCommit
	if err := json.NewDecoder(resp.Body).Decode(&ghCommits); err != nil {
		return nil, err
	}

	out := make([]Commit, 0, len(ghCommits))
	for _, gc := range ghCommits {
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
