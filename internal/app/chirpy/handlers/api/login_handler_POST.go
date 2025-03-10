package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/auth"
	"github.com/johndosdos/chirpy/internal/database"
)

func Login(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type response struct {
			ID           uuid.UUID `json:"id"`
			CreatedAt    time.Time `json:"created_at"`
			UpdatedAt    time.Time `json:"updated_at"`
			Email        string    `json:"email"`
			Token        string    `json:"token"`
			RefreshToken string    `json:"refresh_token"`
			IsChirpyRed  bool      `json:"is_chirpy_red"`
		}

		var req request

		// decode request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Invalid request: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
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
		//
		// access token expire after 1 hour
		jwt, err := auth.MakeJWT(user.ID, cfg.Secret, time.Duration(1)*time.Hour)
		if err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// generate a refresh token and return the token together with
		// the access token (JWT)
		//
		// save refresh token to DB
		//
		// refresh token expire after 60 days
		newRefreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		refreshToken, err := cfg.DB.MakeRefreshToken(r.Context(), database.MakeRefreshTokenParams{
			Token:     newRefreshToken,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(60 * 24 * time.Hour),
		})
		if err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// encode the response
		if err := json.NewEncoder(w).Encode(response{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        jwt,
			RefreshToken: refreshToken.Token,
			IsChirpyRed:  user.IsChirpyRed,
		}); err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
