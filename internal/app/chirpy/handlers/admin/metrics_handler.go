package admin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
)

func GetHits(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		serverHits := cfg.FileserverHits.Load()
		fmt.Fprintf(w, `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Document</title>
    </head>
    <body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
    </body>
</html>
		`, serverHits)
	})
}

func ResetMetrics(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		if cfg.Platform != "dev" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)

		cfg.FileserverHits.Store(0)
		serverHits := cfg.FileserverHits.Load()
		fmt.Fprintf(w, "Hits: %d\n", serverHits)

		// reset users table
		resetUsersDb(cfg, w, r)
	})
}

func resetUsersDb(cfg *chirpy.ApiConfig, w http.ResponseWriter, r *http.Request) {
	err := cfg.DB.DeleteUsers(r.Context())
	if err != nil {
		log.Println("failed to reset users table: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
