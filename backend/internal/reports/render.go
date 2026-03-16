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

func (r *ReleaseDiffReport) ToMarkdown() string {
	var b strings.Builder

	// Título
	fmt.Fprintf(&b, "# Release diff `%s...%s`\n\n", r.FromTag, r.ToTag)
	fmt.Fprintf(&b, "- Repositório: `%s/%s`\n", r.Repository.Owner, r.Repository.Name)
	fmt.Fprintf(&b, "- De (from): `%s`\n", r.FromTag)
	fmt.Fprintf(&b, "- Para (to): `%s`\n", r.ToTag)
	fmt.Fprintf(&b, "- Gerado em: %s\n\n", formatTimeVal(r.GeneratedAt))

	// Resumo
	fmt.Fprintf(&b, "## Resumo\n\n")
	fmt.Fprintf(&b, "- Total de PRs: **%d**\n", r.Summary.TotalPRs)

	if len(r.Summary.ByType) > 0 {
		fmt.Fprintf(&b, "- PRs por tipo:\n")
		// ordenar tipos para saída estável
		type typeCount struct {
			typ   ReleasePRType
			count int
		}
		var items []typeCount
		for typ, count := range r.Summary.ByType {
			items = append(items, typeCount{typ: typ, count: count})
		}
		sort.Slice(items, func(i, j int) bool {
			return string(items[i].typ) < string(items[j].typ)
		})
		for _, it := range items {
			fmt.Fprintf(&b, "  - `%s`: **%d**\n", it.typ, it.count)
		}
	}
	fmt.Fprintf(&b, "\n")

	// Ordenar PRs por data de merge desc (ou número, se não tiver merged_at)
	prs := make([]ReleasePR, len(r.PRs))
	copy(prs, r.PRs)
	sort.Slice(prs, func(i, j int) bool {
		ti := time.Time{}
		if prs[i].MergedAt != nil {
			ti = *prs[i].MergedAt
		}
		tj := time.Time{}
		if prs[j].MergedAt != nil {
			tj = *prs[j].MergedAt
		}
		// mais recente primeiro
		return ti.After(tj)
	})

	fmt.Fprintf(&b, "## PRs incluídas entre `%s` e `%s`\n\n", r.FromTag, r.ToTag)

	for _, pr := range prs {
		fmt.Fprintf(&b, "### PR #%d - %s\n\n", pr.Number, pr.Title)
		fmt.Fprintf(&b, "- Link: %s\n", pr.URL)

		if pr.AuthorLogin != "" {
			fmt.Fprintf(&b, "- Autor: `%s`\n", pr.AuthorLogin)
		}
		if pr.MergedAt != nil {
			fmt.Fprintf(&b, "- Mergeado em: %s\n", formatTime(pr.MergedAt))
		}
		if len(pr.Labels) > 0 {
			fmt.Fprintf(&b, "- Labels: ")
			labels := make([]string, len(pr.Labels))
			for i, l := range pr.Labels {
				labels[i] = fmt.Sprintf("`%s`", l)
			}
			fmt.Fprintf(&b, "%s\n", strings.Join(labels, ", "))
		}
		if pr.TypeClassification != "" && pr.TypeClassification != ReleasePRTypeUnknown {
			fmt.Fprintf(&b, "- Tipo: `%s`\n", pr.TypeClassification)
		}

		fmt.Fprintf(&b, "\n")
	}

	return b.String()
}
