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

func UpdateUserInfo(cfg *chirpy.ApiConfig) http.Handler {
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

		// retrieve user token from authorization header and
		// then we validate it
		tokenString, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Println("invalid authorization header: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		userID, err := auth.ValidateJWT(tokenString, cfg.Secret)
		if err != nil {
			log.Println("invalid access token: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// decode client request; email and password in this case
		//
		// and then we hash user password and update user's database entry
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println("failed to decode request body: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Println("failed to hash user password: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
			Email:          req.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		})
		if err != nil {
			log.Println("failed to update user info in the database: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// return 200 OK and response struct
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(response{
			Id:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		})
		if err != nil {
			log.Println("failed to encode server response: ", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})
}
