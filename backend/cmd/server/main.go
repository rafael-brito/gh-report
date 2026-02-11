package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rafael-brito/gh-report/backend/internal/api"
	"github.com/rafael-brito/gh-report/backend/internal/githubclient"
	"github.com/rafael-brito/gh-report/backend/internal/reports"
)

func main() {
	// Em produção, leia de config/env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Aqui futuramente você vai criar um http.Client com transporte customizado
	ghToken := os.Getenv("GITHUB_TOKEN") // provisório p/ testes locais
	if ghToken == "" {
		log.Println("WARNING: GITHUB_TOKEN não definido, chamadas ao GitHub irão falhar")
	}

	ghClient := githubclient.NewHTTPClient(ghToken)

	fileHistorySvc := reports.NewFileHistoryService(ghClient)
	releaseDiffSvc := reports.NewReleaseDiffService(ghClient)

	router := api.NewRouter(fileHistorySvc, releaseDiffSvc)

	log.Printf("Servidor ouvindo em :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
