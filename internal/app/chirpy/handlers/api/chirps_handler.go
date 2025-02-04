package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/johndosdos/chirpy/internal/app/chirpy"
)

func ProcessChirp(cfg *chirpy.ApiConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// encode a go struct to json using json.Marshal.

		// we return a response after encoding when the data being sent is
		// valid or invalid; true or false and http status codes.

		type response struct {
			Id         uuid.UUID `json:"id"`
			Created_at time.Time `json:"created_at"`
			Updated_at time.Time `json:"updated_at"`
			UserId     uuid.UUID `json:"user_id"`
			Body       string    `json:"body"`
			Error      string    `json:"error"`
		}

		type request struct {
			Body   string    `json:"body"`
			UserId uuid.UUID `json:"user_id"`
		}

		const MAX_CHAR_LEN = 140
		var req request

		// first, decode request body
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			log.Print(err)
		}

		// then, sanitize the request body by checking for profanity.
		sanitizedBody := sanitizeBody(req.Body)
		req.Body = sanitizedBody

		// then, return a 400 http error (bad request) if char > 140. else,
		// return 200 (OK)
		w.Header().Set("Content-Type", "application/json")

		// AVOID MAGIC NUMBERS, i.e., MAX_CHAR_LEN
		if len(sanitizedReq) <= MAX_CHAR_LEN {
			w.WriteHeader(http.StatusOK)

			encoder := json.NewEncoder(w)
			err := encoder.Encode(response{
				Cleaned_body: sanitizedReq,
			})
			if err != nil {
				log.Println("Failed to encode response json: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)

			encoder := json.NewEncoder(w)
			err := encoder.Encode(response{
				Error: "Chirp is too long. Max character length is 140.",
				Body:  sanitizedBody,
			})
			if err != nil {
				log.Println("Chirp is more than 140 chars: ", err)
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
		}

	})
}

func sanitizeBody(body string) string {
	// the reason for using structs as values is because nil structs
	// don't allocate memory. bool as values, on the other hand, will
	// allocate memory. we only need this map for existence checks.
	// also maps have fast lookups.

	// profane words with punctuations are not sanitized.
	profanityMap := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	splitBody := strings.Fields(body)
	for pos, word := range splitBody {
		loword := strings.ToLower(word)
		if _, ok := profanityMap[loword]; ok {
			splitBody[pos] = "****"
		}
	}

	return strings.Join(splitBody, " ")
}
