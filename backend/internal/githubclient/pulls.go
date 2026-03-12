package githubclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ghPR struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	HTMLURL string `json:"html_url"`
	User    struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"user"`
	MergedAt *time.Time `json:"merged_at"`
	Labels   []struct {
		Name string `json:"name"`
	} `json:"labels"`
}

func (c *httpClient) ListPRsByCommit(ctx context.Context, owner, repo, sha string) ([]PRShort, error) {
	cacheKey := fmt.Sprintf("%s/%s/%s", owner, repo, sha)

	c.mu.Lock()
	if prs, ok := c.prsByCommitCache[cacheKey]; ok {
		c.mu.Unlock()
		return prs, nil
	}
	c.mu.Unlock()

	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s/pulls", c.baseURL, owner, repo, sha)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	// Header recomendado para esse endpoint
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Sem PRs associados a esse commit (normal)
		return []PRShort{}, nil
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github ListPRsByCommit: status %d", resp.StatusCode)
	}

	var ghPRs []ghPR
	if err := json.NewDecoder(resp.Body).Decode(&ghPRs); err != nil {
		return nil, err
	}

	out := make([]PRShort, 0, len(ghPRs))
	for _, p := range ghPRs {
		labels := make([]string, 0, len(p.Labels))
		for _, l := range p.Labels {
			labels = append(labels, l.Name)
		}
		out = append(out, PRShort{
			Number:          p.Number,
			Title:           p.Title,
			URL:             p.HTMLURL,
			AuthorLogin:     p.User.Login,
			AuthorAvatarURL: p.User.AvatarURL,
			MergedAt:        p.MergedAt,
			Labels:          labels,
		})
	}

	c.mu.Lock()
	c.prsByCommitCache[cacheKey] = out
	c.mu.Unlock()

	return out, nil
}
