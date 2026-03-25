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

// Interfaces usadas principalmente para testes/mocks se você quiser no futuro.
type FileHistoryHandlerDeps interface {
	GetFileHistoryReport(ctx context.Context, params reports.FileHistoryParams) (*reports.FileHistoryReport, error)
}

type ReleaseDiffHandlerDeps interface {
	GetReleaseDiffReport(ctx context.Context, params reports.ReleaseDiffParams) (*reports.ReleaseDiffReport, error)
}

// ReportsHandler agora é "com auth": escolhe o token por request e cria o client/services.
type ReportsHandler struct {
	tokenProvider TokenProvider
	clientFactory GitHubClientFactory
}

func NewReportsHandler(tokenProvider TokenProvider, clientFactory GitHubClientFactory) *ReportsHandler {
	return &ReportsHandler{
		tokenProvider: tokenProvider,
		clientFactory: clientFactory,
	}
}

const dummyUserID = "dev-local"

// Cria services para esta request usando o token apropriado.
func (h *ReportsHandler) newServicesForRequest(r *http.Request) (reports.FileHistoryService, reports.ReleaseDiffService) {
	token := h.tokenProvider.TokenForRequest(r)
	client := h.clientFactory.ForToken(token)

	fileHistorySvc := reports.NewFileHistoryService(client)
	releaseDiffSvc := reports.NewReleaseDiffService(client)

	return fileHistorySvc, releaseDiffSvc
}

func parseRepoParam(repoStr string) (reports.RepositoryRef, error) {
	parts := strings.Split(repoStr, "/")
	if len(parts) != 2 {
		return reports.RepositoryRef{}, fmt.Errorf("invalid repo format, expected owner/repo")
	}
	return reports.RepositoryRef{Owner: parts[0], Name: parts[1]}, nil
}

func (h *ReportsHandler) HandleFileHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := dummyUserID

	fileHistorySvc, _ := h.newServicesForRequest(r)

	repoStr := r.URL.Query().Get("repo")
	file := r.URL.Query().Get("file")
	limitStr := r.URL.Query().Get("limit")
	modeStr := r.URL.Query().Get("mode")
	format := strings.ToLower(r.URL.Query().Get("format"))

	if format == "" {
		format = "json"
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

	report, err := fileHistorySvc.GetFileHistoryReport(ctx, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch format {
	case "markdown", "md":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		md := report.ToMarkdown()
		_, _ = w.Write([]byte(md))
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		csvStr, err := report.ToCSV()
		if err != nil {
			http.Error(w, "failed to render csv: "+err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(csvStr))
	default:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(report)
	}
}

func (h *ReportsHandler) HandleReleaseDiff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := dummyUserID

	_, releaseDiffSvc := h.newServicesForRequest(r)

	repoStr := r.URL.Query().Get("repo")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	format := strings.ToLower(r.URL.Query().Get("format"))

	if format == "" {
		format = "json"
	}

	repoRef, err := parseRepoParam(repoStr)
	if err != nil || from == "" || to == "" {
		http.Error(w, "invalid repo, from or to", http.StatusBadRequest)
		return
	}

	params := reports.ReleaseDiffParams{
		Repo:   repoRef,
		From:   from,
		To:     to,
		UserID: userID,
	}

	report, err := releaseDiffSvc.GetReleaseDiffReport(ctx, params)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "base or head not found") {
			http.Error(w, "invalid from/to ref or tag not found", http.StatusBadRequest)
			return
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	switch format {
	case "markdown", "md":
		w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		md := report.ToMarkdown()
		_, _ = w.Write([]byte(md))
	case "csv":
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		csvStr, err := report.ToCSV()
		if err != nil {
			http.Error(w, "failed to render csv: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		_, _ = w.Write([]byte(csvStr))
	default:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(report)
	}
}
