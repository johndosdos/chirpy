package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/auth"
	"github.com/johndosdos/chirpy/internal/database"
)

func Revoke(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get refresh token from client request
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Println("failed to extract Bearer token: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// retrieve refresh token user from DB and update user's refresh token
		// by updating the DB
		user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), tokenString)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("user not found: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				log.Println("failed to get refresh token user: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		err = cfg.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
			RevokedAt: sql.NullTime{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			UpdatedAt: time.Now().UTC(),
			UserID:    user.UserID,
		})
		if err != nil {
			log.Println("failed to update user refresh token: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
