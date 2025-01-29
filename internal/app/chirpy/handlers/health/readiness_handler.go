package health

import (
	"log"
	"net/http"
)

func Check(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		/*
			http.ResponseWriter handles the response that our server sends back
			to the client.
		*/

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		/*
			technically, the error won't be returned to the client after w.Write().
			it's just here for logging purposes.

			it's a good practice tho.
		*/
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Printf("failed to write healthz response: %v", err)
		}
	})
}
