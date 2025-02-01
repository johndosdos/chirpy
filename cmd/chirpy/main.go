package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/admin/health"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/admin/metric"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/api"
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

	mux.Handle("POST /api/validate_chirp", api.ValidateChirp())

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server starting at port 8080...")
	log.Fatal(server.ListenAndServe())
}
