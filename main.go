package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/admin/health"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/admin/metric"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/api"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/api/users"
	"github.com/johndosdos/chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// DATABASE INIT...

	// Load .env file from project root. Our .env file contain sensitive
	// keys and information. Please add to .gitignore
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file: ", err)
	}

	dbUrl := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("failed to initialize db: ", err)
	}

	dbQueries := database.New(db)

	// SERVER INIT...
	mux := http.NewServeMux()
	apiCfg := &chirpy.ApiConfig{
		DB:       dbQueries,
		Platform: platform,
	}

	// check file server readiness.
	health.Check(mux)

	// strip the prefix "/app/" from the URL path for proper routing.
	// URL path != file path on the server.
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("web/")))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(fileServer))

	mux.Handle("GET /admin/metrics", metric.GetHits(apiCfg))
	mux.Handle("POST /admin/reset", metric.ResetMetrics(apiCfg))

	mux.Handle("POST /api/validate_chirp", api.ValidateChirp())
	mux.Handle("POST /api/users", users.CreateUser(apiCfg))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server starting at port 8080...")
	log.Fatal(server.ListenAndServe())
}
