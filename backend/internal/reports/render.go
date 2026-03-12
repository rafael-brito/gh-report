package reports

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format("2006-01-02 15:04:05 MST")
}

func formatTimeVal(t time.Time) string {
	return t.UTC().Format("2006-01-02 15:04:05 MST")
}

func (r *FileHistoryReport) ToMarkdown() string {
	var b strings.Builder

	// Título
	fmt.Fprintf(&b, "# Histórico do arquivo `%s`\n\n", r.FilePath)
	fmt.Fprintf(&b, "- Repositório: `%s/%s`\n", r.Repository.Owner, r.Repository.Name)
	fmt.Fprintf(&b, "- Modo: **%s**\n", r.Mode)
	fmt.Fprintf(&b, "- Limite: %d\n", r.Limit)
	fmt.Fprintf(&b, "- Gerado em: %s\n\n", formatTimeVal(r.GeneratedAt))

	// Stats
	fmt.Fprintf(&b, "## Estatísticas\n\n")
	fmt.Fprintf(&b, "- Total de commits considerados: **%d**\n", r.Stats.TotalCommits)
	fmt.Fprintf(&b, "- Total de PRs: **%d**\n\n", r.Stats.TotalPRs)

	if len(r.Stats.TopAuthors) > 0 {
		fmt.Fprintf(&b, "### Top autores (por commits)\n\n")
		fmt.Fprintf(&b, "| Autor | Commits |\n")
		fmt.Fprintf(&b, "|-------|---------|\n")
		for _, a := range r.Stats.TopAuthors {
			fmt.Fprintf(&b, "| `%s` | %d |\n", a.Login, a.Commits)
		}
		fmt.Fprintf(&b, "\n")
	}

	// Ordenar entries por OrderTs desc (mais recente primeiro)
	entries := make([]FileHistoryEntry, len(r.Entries))
	copy(entries, r.Entries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].OrderTs.After(entries[j].OrderTs)
	})

	switch r.Mode {
	case FileHistoryModePRs:
		b.WriteString("## PRs que modificaram o arquivo\n\n")
		for _, e := range entries {
			if e.Type != FileHistoryEntryTypePR || e.PR == nil {
				continue
			}
			pr := e.PR
			fmt.Fprintf(&b, "### PR #%d - %s\n\n", pr.PRNumber, pr.PRTitle)
			fmt.Fprintf(&b, "- Link: %s\n", pr.PRURL)
			if pr.AuthorLogin != "" {
				fmt.Fprintf(&b, "- Autor: `%s`\n", pr.AuthorLogin)
			}
			if pr.PRMergedAt != nil {
				fmt.Fprintf(&b, "- Mergeado em: %s\n", formatTime(pr.PRMergedAt))
			}
			if len(pr.Commits) > 0 {
				fmt.Fprintf(&b, "- Commits que tocaram o arquivo: %d\n\n", len(pr.Commits))
				fmt.Fprintf(&b, "| SHA | Autor | Data | Mensagem |\n")
				fmt.Fprintf(&b, "|-----|-------|------|----------|\n")
				for _, c := range pr.Commits {
					shortSHA := c.SHA
					if len(shortSHA) > 8 {
						shortSHA = shortSHA[:8]
					}
					msg := c.Message
					msg = strings.ReplaceAll(msg, "\n", " ")
					if len(msg) > 80 {
						msg = msg[:77] + "..."
					}
					fmt.Fprintf(
						&b,
						"| [`%s`](%s) | `%s` | %s | %s |\n",
						shortSHA,
						c.URL,
						nullOr(c.AuthorLogin, "-"),
						formatTimeVal(c.CommittedAt),
						escapePipes(msg),
					)
				}
				fmt.Fprintf(&b, "\n")
			} else {
				fmt.Fprintf(&b, "- Nenhum commit listado para este PR.\n\n")
			}
		}
	default:
		// modo commits
		b.WriteString("## Commits que modificaram o arquivo\n\n")
		fmt.Fprintf(&b, "| SHA | Autor | Data | Mensagem |\n")
		fmt.Fprintf(&b, "|-----|-------|------|----------|\n")
		for _, e := range entries {
			if e.Type != FileHistoryEntryTypeCommit || e.Commit == nil {
				continue
			}
			c := e.Commit
			shortSHA := c.SHA
			if len(shortSHA) > 8 {
				shortSHA = shortSHA[:8]
			}
			msg := c.Message
			msg = strings.ReplaceAll(msg, "\n", " ")
			if len(msg) > 80 {
				msg = msg[:77] + "..."
			}
			fmt.Fprintf(
				&b,
				"| [`%s`](%s) | `%s` | %s | %s |\n",
				shortSHA,
				c.URL,
				nullOr(c.AuthorLogin, "-"),
				formatTimeVal(c.CommittedAt),
				escapePipes(msg),
			)
		}
		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}

func nullOr(v, alt string) string {
	if v == "" {
		return alt
	}
	return v
}

func escapePipes(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
