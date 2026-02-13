package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/rafael-brito/gh-report/backend/internal/reports"
)

type FileHistoryHandlerDeps interface {
	GetFileHistoryReport(ctx context.Context, params reports.FileHistoryParams) (*reports.FileHistoryReport, error)
}

type ReleaseDiffHandlerDeps interface {
	GetReleaseDiffReport(ctx context.Context, params reports.ReleaseDiffParams) (*reports.ReleaseDiffReport, error)
}

type ReportsHandler struct {
	fileHistoryService reports.FileHistoryService
	releaseDiffService reports.ReleaseDiffService
}

func NewReportsHandler(fh reports.FileHistoryService, rd reports.ReleaseDiffService) *ReportsHandler {
	return &ReportsHandler{
		fileHistoryService: fh,
		releaseDiffService: rd,
	}
}

func parseRepoParam(repoStr string) (reports.RepositoryRef, error) {
	parts := strings.Split(repoStr, "/")
	if len(parts) != 2 {
		return reports.RepositoryRef{}, fmt.Errorf("invalid repo format, expected owner/repo")
	}
	return reports.RepositoryRef{Owner: parts[0], Name: parts[1]}, nil
}

const dummyUserID = "dev-local"

func (h *ReportsHandler) HandleFileHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := dummyUserID

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

	report, err := h.fileHistoryService.GetFileHistoryReport(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(report)
	default:
		// temporariamente, só JSON para conseguir testar rápido
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(report)
	}
}

func (h *ReportsHandler) HandleReleaseDiff(w http.ResponseWriter, r *http.Request) {
	// Placeholder: ainda não implementado
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
