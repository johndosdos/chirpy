package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/health"
)

func main() {
	mux := http.NewServeMux()

	// check file server readiness.
	health.Check(mux)

	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("web/")))
	mux.Handle("/app/", fileServer)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server starting at port 8080...")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
