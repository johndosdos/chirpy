package api

import (
	"encoding/json"
	"log"
	"net/http"
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
