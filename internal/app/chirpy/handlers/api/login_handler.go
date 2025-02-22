package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/auth"
)

func Login(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
		}

		type response struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
			Token     string    `json:"token"`
		}

		var req request

		// default expiration time if request ExpiresInSeconds is nil
		//
		// defaults to 3600 seconds or 1 hour
		const DEFAULT_EXPIRATION int = 3600

		// decode request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Invalid request: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// check request body if expires_in_seconds is nil or not
		if req.ExpiresInSeconds <= 0 || req.ExpiresInSeconds > 3600 {
			req.ExpiresInSeconds = DEFAULT_EXPIRATION
		}

		// get user info by email
		user, err := cfg.DB.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		// compare request password to the stored, hashed password
		if err := auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		// generate JWT
		//
		// note that we need to multipy time.Duration by time.Second since
		// time.Duration will convert to time in nanoseconds
		jwt, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(req.ExpiresInSeconds)*time.Second)
		if err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// encode the response
		if err := json.NewEncoder(w).Encode(response{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     jwt,
		}); err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
