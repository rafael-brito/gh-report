package reports

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strings"
)

func (r *FileHistoryReport) ToCSV() (string, error) {
	var sb strings.Builder
	w := csv.NewWriter(&sb)

	// Cabeçalho
	if r.Mode == FileHistoryModePRs {
		if err := w.Write([]string{
			"repo_owner",
			"repo_name",
			"file_path",
			"pr_number",
			"pr_title",
			"pr_url",
			"pr_merged_at",
			"commit_sha",
			"commit_author",
			"commit_date",
			"commit_message",
		}); err != nil {
			return "", err
		}
	} else {
		if err := w.Write([]string{
			"repo_owner",
			"repo_name",
			"file_path",
			"commit_sha",
			"commit_author",
			"commit_date",
			"commit_message",
		}); err != nil {
			return "", err
		}
	}

	entries := make([]FileHistoryEntry, len(r.Entries))
	copy(entries, r.Entries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].OrderTs.After(entries[j].OrderTs)
	})

	switch r.Mode {
	case FileHistoryModePRs:
		for _, e := range entries {
			if e.Type != FileHistoryEntryTypePR || e.PR == nil {
				continue
			}
			pr := e.PR
			for _, c := range pr.Commits {
				row := []string{
					r.Repository.Owner,
					r.Repository.Name,
					r.FilePath,
					fmt.Sprintf("%d", pr.PRNumber),
					pr.PRTitle,
					pr.PRURL,
					formatTime(pr.PRMergedAt),
					c.SHA,
					c.AuthorLogin,
					formatTimeVal(c.CommittedAt),
					c.Message,
				}
				if err := w.Write(row); err != nil {
					return "", err
				}
			}
		}
	default:
		for _, e := range entries {
			if e.Type != FileHistoryEntryTypeCommit || e.Commit == nil {
				continue
			}
			c := e.Commit
			row := []string{
				r.Repository.Owner,
				r.Repository.Name,
				r.FilePath,
				c.SHA,
				c.AuthorLogin,
				formatTimeVal(c.CommittedAt),
				c.Message,
			}
			if err := w.Write(row); err != nil {
				return "", err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return sb.String(), nil
}
