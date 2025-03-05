package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/johndosdos/chirpy/internal/app/chirpy"
)

func WebhookHandler(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// this is the structure for the webhook request
		type request struct {
			Event string `json:"event"`
			Data  struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Println("failed to decode request body: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// return 204 No Content if "event" is anything other than "user.upgraded".
		// we only care if the user has been upgraded to the (hypothetically)
		// premium plan.
		//
		// if "event" is "user.upgraded", it should update the user in the database,
		// and mark that they are a member.
		if req.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// convert userID string to a UUID
		userID, err := uuid.Parse(req.Data.UserID)
		if err != nil {
			log.Println("invalid userID: ", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// we need to check if user exists in the database before we can
		// upgrade them
		_, err = cfg.DB.UpgradeUser(r.Context(), userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println("failed to upgrade user: user not found: ", err)
				http.Error(w, "Not found", http.StatusNotFound)
			} else {
				log.Println("database error: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// if the user is upgraded successfully, the endpoint should
		// respond with a 204 status code and an empty response body.
		w.WriteHeader(http.StatusNoContent)
	})
}
