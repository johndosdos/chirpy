package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/health"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/metric"
)

func main() {
	mux := http.NewServeMux()
	apiCfg := &chirpy.ApiConfig{}

	// check file server readiness.
	health.Check(mux)

	// strip the prefix "/app/" from the URL path for proper routing.
	// URL path != file path on the server.
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("web/")))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(fileServer))
	mux.Handle("GET /admin/metrics", metric.GetHits(apiCfg))
	mux.Handle("POST /admin/reset", metric.ResetMetrics(apiCfg))

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
