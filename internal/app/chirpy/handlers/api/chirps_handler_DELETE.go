package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/auth"
)

func DeleteChirp(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract endpoint id or pathvalue (idk what's the correct term)
		// from URL.
		chirpID, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			log.Println("invalid chirp ID: ", err)
			http.Error(w, "Bad request: invalid chirp ID format", http.StatusBadRequest)
			return
		}

		// authenticate access token and then validate it.
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Println("invalid authorization header: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(tokenString, cfg.Secret)
		if err != nil {
			log.Println("invalid access token: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check if user is the author of the chirp
		chirp, err := cfg.DB.GetChirp(r.Context(), chirpID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("chirp not found: ", err)
				http.Error(w, "Not found: chirp not found", http.StatusNotFound)
			} else {
				log.Println("database error: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		if userID != chirp.UserID {
			log.Println("chirp deletion not allowed.")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// if no errors were present, authorize current user to delete
		// chirp by their id.
		err = cfg.DB.DeleteChirp(r.Context(), chirpID)
		if err != nil {
			log.Println("failed to delete chirp: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// return No Content 204 confirmation on successful deletion
		w.WriteHeader(http.StatusNoContent)
	})
}
