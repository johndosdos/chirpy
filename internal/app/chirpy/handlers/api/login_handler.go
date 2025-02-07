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
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type response struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
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

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// encode the response
		if err := json.NewEncoder(w).Encode(response{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}); err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
