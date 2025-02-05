package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
)

func GetChirps(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get data from db
		chirps, err := cfg.DB.GetChirps(r.Context())
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// encode []chirps directly to w
		encoder := json.NewEncoder(w)
		err = encoder.Encode(chirps)
		if err != nil {
			log.Println("Failed to encode response json: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
