package api

import (
	"net/http"
)

func NewRouter(fileHistoryHandler FileHistoryHandlerDeps, releaseDiffHandler ReleaseDiffHandlerDeps) http.Handler {
	mux := http.NewServeMux()

	reportsHandler := NewReportsHandler(fileHistoryHandler, releaseDiffHandler)

	// Rotas de relatórios
	mux.HandleFunc("/api/reports/file-history", reportsHandler.HandleFileHistory)
	mux.HandleFunc("/api/reports/release-diff", reportsHandler.HandleReleaseDiff)

	// Aqui depois você adiciona /auth/login, /auth/callback, /auth/me

	// No futuro, você pode envolver mux com middlewares de auth, logging, CORS etc.
	return mux
}
