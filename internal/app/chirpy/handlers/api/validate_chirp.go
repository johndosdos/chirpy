package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func ValidateChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// encode a go struct to json using json.Marshal.

		// we return a response after encoding when the data being sent is
		// valid or invalid; true or false and http status codes.

		type response struct {
			Error string `json:"error"`
			Valid bool   `json:"valid"`
		}

		type request struct {
			Body string `json:"body"`
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
		sanitizedReq := sanitizeBody(req.Body)

		// then, return a 400 http error (bad request) if char > 140. else,
		// return 200 (OK)
		w.Header().Set("Content-Type", "application/json")

		// AVOID MAGIC NUMBERS, i.e., MAX_CHAR_LEN
		if len(req.Body) <= MAX_CHAR_LEN {
			w.WriteHeader(http.StatusOK)

			encoder := json.NewEncoder(w)
			err := encoder.Encode(response{
				Valid: true,
			})
			if err != nil { // check/revise error handling !!
				log.Print(err)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)

			encoder := json.NewEncoder(w)
			err := encoder.Encode(response{
				Error: "Chirp is too long. Max character length is 140.",
				Valid: false,
			})
			if err != nil { // check/revise error handling !!
				log.Print(err)
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
