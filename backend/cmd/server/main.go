package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/rafael-brito/gh-report/backend/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Println("WARNING: GITHUB_TOKEN não definido; se X-GitHub-Token não for enviado, chamadas ao GitHub irão falhar")
	}

	tokenProvider := api.NewSimpleTokenProvider()
	clientFactory := api.NewGitHubClientFactory()

	apiRouter := api.NewRouterWithAuth(tokenProvider, clientFactory)

	// Diretório do build do frontend
	frontendDist := os.Getenv("FRONTEND_DIST_DIR")
	if frontendDist == "" {
		// em desenvolvimento local, após `npm run build` no frontend
		frontendDist = "./frontend/dist"
	}

	log.Printf("Usando FRONTEND_DIST_DIR=%s\n", frontendDist)

	// Mux raiz: combina API + static/frontend
	rootMux := http.NewServeMux()

	// Monta /api/ e /healthz no mux raiz
	rootMux.Handle("/api/", apiRouter)
	rootMux.HandleFunc("/healthz", api.HandleHealth)

	// Static files do Vite (assets)
	fs := http.FileServer(http.Dir(frontendDist))
	rootMux.Handle("/assets/", fs)
	rootMux.Handle("/favicon.ico", fs)
	rootMux.Handle("/gh-report-icon.svg", fs) // opcional, se estiver no dist

	// Catch-all para SPA (React Router)
	rootMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexPath := filepath.Join(frontendDist, "index.html")
		http.ServeFile(w, r, indexPath)
	})

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      rootMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Servidor ouvindo em :%s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
