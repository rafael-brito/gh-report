package api

import (
	"net/http"

	"github.com/rafael-brito/gh-report/backend/internal/reports"
)

func NewRouter(
	fileHistorySvc reports.FileHistoryService,
	releaseDiffSvc reports.ReleaseDiffService,
) http.Handler {
	mux := http.NewServeMux()

	reportsHandler := NewReportsHandler(fileHistorySvc, releaseDiffSvc)

	// Rotas de relatórios
	mux.HandleFunc("/api/reports/file-history", reportsHandler.HandleFileHistory)
	mux.HandleFunc("/api/reports/release-diff", reportsHandler.HandleReleaseDiff)

	// Aqui depois você adiciona /auth/login, /auth/callback, /auth/me

	return mux
}
