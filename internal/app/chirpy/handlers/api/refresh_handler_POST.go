package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/auth"
)

func Refresh(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Token string `json:"token"`
		}

		// get access token from client request
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Println("failed to extract bearer token: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// using the Bearer token, check if token exists in the
		// refresh_tokens SQL table
		user, err := cfg.DB.GetUserFromRefreshToken(r.Context(), token)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("invalid refresh token: not found in database or expired")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			} else {
				log.Println("database error: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// check expiration
		if user.ExpiresAt.Before(time.Now()) {
			log.Println("invalid refresh token: expired")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// check revoke validity
		if user.RevokedAt.Valid {
			log.Println("invalid refresh token: revoked")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// create new access token for user after checking
		tokenString, err := auth.MakeJWT(user.UserID, cfg.Secret, time.Duration(1*time.Hour))
		if err != nil {
			log.Println("failed create JWT: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// we need to return the new access token to the client and return 200 OK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(response{Token: tokenString})
		if err != nil {
			log.Println("failed to encode JSON response: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
