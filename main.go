package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/johndosdos/chirpy/internal/app/chirpy"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/admin"
	"github.com/johndosdos/chirpy/internal/app/chirpy/handlers/api"
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
	secret := os.Getenv("SECRET")

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
		Secret:   secret,
	}

	// check file server readiness.
	admin.Check(mux)

	// strip the prefix "/app/" from the URL path for proper routing.
	// URL path != file path on the server.
	fileServer := http.StripPrefix("/app/", http.FileServer(http.Dir("web/")))

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(fileServer))

	mux.Handle("GET /admin/metrics", admin.GetHits(apiCfg))
	mux.Handle("POST /admin/reset", admin.ResetMetrics(apiCfg))

	mux.Handle("GET /api/chirps/{chirpID}", api.GetChirp(apiCfg))
	mux.Handle("GET /api/chirps", api.GetChirps(apiCfg))
	mux.Handle("POST /api/chirps", api.ProcessChirp(apiCfg))
	mux.Handle("DELETE /api/chirps/{chirpID}", api.DeleteChirp(apiCfg))

	mux.Handle("POST /api/users", api.CreateUser(apiCfg))
	mux.Handle("PUT /api/users", api.UpdateUserInfo(apiCfg))

	mux.Handle("POST /api/login", api.Login(apiCfg))

	mux.Handle("POST /api/refresh", api.Refresh(apiCfg))

	mux.Handle("POST /api/revoke", api.Revoke(apiCfg))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server starting at port 8080...")
	log.Fatal(server.ListenAndServe())
}
