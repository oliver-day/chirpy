package main

import (
	"log"
	"net/http"
)

func main() {
	const filePathRoot = "."
	const port = "8080"

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(filePathRoot))
	mux.Handle("/app/*", http.StripPrefix("/app", fileServer))
	mux.HandleFunc("/healthz", handlerReadiness)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server files from %s on port %s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
