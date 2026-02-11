package api

import (
	"encoding/json"
	"gh-report/internal/reports"
	"net/http"
	"strconv"
	"strings"
)

type ReportsHandler struct {
	FileHistory reports.FileHistoryService
	ReleaseDiff reports.ReleaseDiffService
}

func parseRepoParam(repoStr string) (reports.RepositoryRef, error) {
	parts := strings.Split(repoStr, "/")
	if len(parts) != 2 {
		return reports.RepositoryRef{}, ErrBadRequest // defina um erro apropriado
	}
	return reports.RepositoryRef{Owner: parts[0], Name: parts[1]}, nil
}

func (h *ReportsHandler) HandleFileHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Aqui, suponho que você já tenha middleware que injeta userID
	userID := ctx.Value(UserIDKey).(string)

	repoStr := r.URL.Query().Get("repo")
	file := r.URL.Query().Get("file")
	limitStr := r.URL.Query().Get("limit")
	modeStr := r.URL.Query().Get("mode")
	format := r.URL.Query().Get("format")

	if format == "" {
		format = "markdown"
	}
	if modeStr == "" {
		modeStr = string(reports.FileHistoryModePRs)
	}
	if limitStr == "" {
		limitStr = "10"
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	repoRef, err := parseRepoParam(repoStr)
	if err != nil || file == "" {
		http.Error(w, "invalid repo or file", http.StatusBadRequest)
		return
	}

	params := reports.FileHistoryParams{
		Repo:   repoRef,
		File:   file,
		Limit:  limit,
		Mode:   reports.FileHistoryMode(modeStr),
		UserID: userID,
	}

	report, err := h.FileHistory.GetFileHistoryReport(ctx, params)
	if err != nil {
		// mapear erros específicos (rate limit, forbidden, etc.)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	case "markdown":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		if err := writeFileHistoryMarkdown(w, report); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		if err := writeFileHistoryCSV(w, report); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		http.Error(w, "unsupported format", http.StatusBadRequest)
	}
}
