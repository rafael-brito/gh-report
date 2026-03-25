package api

import "net/http"

func NewRouterWithAuth(
	tokenProvider TokenProvider,
	clientFactory GitHubClientFactory,
) http.Handler {
	mux := http.NewServeMux()

	reportsHandler := NewReportsHandler(tokenProvider, clientFactory)

	mux.HandleFunc("/api/reports/file-history", reportsHandler.HandleFileHistory)
	mux.HandleFunc("/api/reports/release-diff", reportsHandler.HandleReleaseDiff)

	return mux
}
