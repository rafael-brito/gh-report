package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rafael-brito/gh-report/backend/internal/api"
)

func main() {
	// Em produção, leia de config/env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// AVISO se não houver token global (ainda serve como fallback)
	if os.Getenv("GITHUB_TOKEN") == "" {
		log.Println("WARNING: GITHUB_TOKEN não definido; se X-GitHub-Token não for enviado, chamadas ao GitHub irão falhar")
	}

	tokenProvider := api.NewSimpleTokenProvider()
	clientFactory := api.NewGitHubClientFactory()

	router := api.NewRouterWithAuth(tokenProvider, clientFactory)

	log.Printf("Servidor ouvindo em :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
