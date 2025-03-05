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

// CreateUser expects an email json field from the http request.
func CreateUser(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		type response struct {
			Id          uuid.UUID `json:"id"`
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			Email       string    `json:"email"`
			IsChirpyRed bool      `json:"is_chirpy_red"`
		}

		var req request

		// parse and decode request.

		// error returned from decoding is usually caused
		// by the client side, so we return http error 400
		// (bad request).
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			log.Println("Invalid JSON request: ", err)
			http.Error(w, "Bad client request", http.StatusBadRequest)
			return
		}

		// hash user password before storing to database
		hashedPw, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Println("Failed to hash password: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// return http error 500 since the error is usually caused
		// by the server.
		user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPw,
		})
		if err != nil {
			log.Println("Unexpected error: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// encode and return the response to the client.
		// return http status 201 (Created)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// encoding error is usually a server-side problem, hence
		// http error 500.
		encoder := json.NewEncoder(w)
		err = encoder.Encode(response{
			Id:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		})
		if err != nil {
			log.Println("failed to encode response: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
