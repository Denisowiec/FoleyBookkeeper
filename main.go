package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Denisowiec/FoleyBookkeeper/internal/db"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db_url string
	db     db.Queries
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
	dbase, err := sql.Open("postgres", cfg.db_url)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	dbQueries := db.New(dbase)
	cfg.db = *dbQueries

	mux := http.NewServeMux()

}
