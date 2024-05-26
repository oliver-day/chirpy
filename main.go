package main

import (
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filePathRoot = "."
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: 0,
	}
	mux := http.NewServeMux()

	fileServer := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	mux.Handle("/app/*", fileServer)

	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("GET /reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server files from %s on port %s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
