package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/database"
)

func GetChirp(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// parse path wildcard (here, it is the user id in /api/chirps/{userID})
		// using http.Request.PathValue.
		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			log.Println("invalid chirp ID: ", err)
			http.Error(w, "Bad request: invalid chirp ID format", http.StatusBadRequest)
			return
		}

		// retrieve user chirp from database.
		//
		// because 'emit_json_tags' option in sqlc is set to true, we don't need
		// explicity create a response struct with json tags. check chirps_handler_POST.go
		// for examples. idk if this is a good practice.
		chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "User chirp data not found", http.StatusNotFound)
			return
		}

		// write to w, send response.
		//
		// WriteHeader will be implicitly called, with 200 OK, at the first
		// successful call to w.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(chirp)
		if err != nil {
			log.Println("Failed to encode response json: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}

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

		// extract optional query parameter 'author_id' from URL
		authorID := r.URL.Query().Get("author_id")

		// return the array of chirps if author_id is nil
		if authorID == "" {
			// encode []chirps directly to w
			encoder := json.NewEncoder(w)
			err = encoder.Encode(chirps)
			if err != nil {
				log.Println("Failed to encode response json: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// return single chirp if author_id is not nil
		//
		// in case of multiple chirps with the same author_id, we
		// need to store them in a slice. This prevents multiple
		// writes to w, which is efficient
		var filteredChirps []database.Chirp

		for _, chirp := range chirps {
			if chirp.UserID != uuid.Nil && chirp.UserID.String() == authorID {
				filteredChirps = append(filteredChirps, chirp)
			}
		}

		encoder := json.NewEncoder(w)
		err = encoder.Encode(filteredChirps)
		if err != nil {
			log.Println("Failed to encode response json: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
