package metric

import (
	"fmt"
	"net/http"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
)

func GetHits(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		serverHits := cfg.FileserverHits.Load()
		fmt.Fprintf(w, "Hits: %d\n", serverHits)
	})
}

func ResetMetrics(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Store(0)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		serverHits := cfg.FileserverHits.Load()
		fmt.Fprintf(w, "Hits: %d\n", serverHits)
	})
}
