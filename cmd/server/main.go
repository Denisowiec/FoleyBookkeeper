package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db_url                 string
	db                     db.Queries
	secret                 string
	jwtExpirationTime      time.Duration
	refTokenExpirationTime time.Duration
	listen_port            string
}

func main() {
	fmt.Println("Welcome to FoleyBookkeeper!")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading configuration from .env file: ", err)
	}

	// cfg is the apiConfig instance that the http server will operate on
	var cfg apiConfig
	cfg.db_url = os.Getenv("DB_URL")
	cfg.secret = os.Getenv("SECRET_KEY")
	cfg.listen_port = os.Getenv("SERVER_LISTEN_PORT")

	// JWT expiration time is provided in .env file as number of seconds
	// It gets converted to time.Duration
	jwtExpirationSeconds, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_TIME"))
	if err != nil {
		log.Fatal("Error processing JWT_EXPIRATION_TIME env variable:", err)
	}
	cfg.jwtExpirationTime = time.Duration(jwtExpirationSeconds) * time.Second

	// Refresh tokens expiration time is provided in .env file as number of
	// seconds. It gets converted to time.Duration
	refTokenExpirationSeconds, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION_TIME"))
	if err != nil {
		log.Fatal("Error processing REFRESH_TOKEN_EXPIRATION_TIME env variable:", err)
	}
	cfg.refTokenExpirationTime = time.Duration(refTokenExpirationSeconds) * time.Second

	// cfg also contains an pointer to the database queries
	dbase, err := sql.Open("postgres", cfg.db_url)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	dbQueries := db.New(dbase)
	cfg.db = *dbQueries

	// Here the api handlers are set up
	mux := http.NewServeMux()

	// User and login related
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefreshToken)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	//mux.HandleFunc("PUT /api/users/{userid}", cfg.handlerUpdateUser)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUserSelf)
	mux.HandleFunc("GET /api/users/{userid}", cfg.handlerGetUser)
	mux.HandleFunc("GET /api/users", cfg.handlerGetUsers)

	// Client related
	mux.HandleFunc("POST /api/clients", cfg.handlerCreateClient)
	mux.HandleFunc("PUT /api/clients/{clientid}", cfg.handlerUpdateClient)
	mux.HandleFunc("GET /api/clients/{clientid}", cfg.handlerGetClientByID)
	mux.HandleFunc("DELETE /api/clients/{clientid}", cfg.handlerDeleteClient)
	mux.HandleFunc("GET /api/clients", cfg.handlerGetClientByName)

	// Project related
	mux.HandleFunc("POST /api/projects", cfg.handlerCreateProject)
	mux.HandleFunc("PUT /api/projects/{projectid}", cfg.handlerUpdateProject)
	mux.HandleFunc("GET /api/projects/{projectid}", cfg.handlerGetProjectByID)
	mux.HandleFunc("DELETE /api/projects/{projectid}", cfg.handlerDeleteProject)
	mux.HandleFunc("GET /api/projects", cfg.handlerGetProjectByTitle)

	// Episode related
	mux.HandleFunc("POST /api/episodes", cfg.handlerCreateEpisode)
	mux.HandleFunc("PUT /api/episodes/{episodeid}", cfg.handlerUpdateEpisode)
	mux.HandleFunc("GET /api/episodes/{episodeid}", cfg.handlerGetEpisodeByID)
	mux.HandleFunc("DELETE /api/episodes/{episodeid}", cfg.handlerDeleteEpisode)
	mux.HandleFunc("GET /api/episodes", cfg.handlerGetEpisodesForProject)

	// Here we create the server
	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.listen_port),
		Handler: mux,
	}

	defer s.Shutdown(context.Background())

	log.Fatal(s.ListenAndServe())
}
