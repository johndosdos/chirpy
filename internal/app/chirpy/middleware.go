package chirpy

import (
	"net/http"
	"sync/atomic"

	"github.com/johndosdos/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
}

// incerment fileserverHits counter everytime a client visits the server,
// especially when the client URL path is "/app/".
func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
