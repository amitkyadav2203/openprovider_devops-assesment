package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"net/url"
	"simplebank/pkg/api"
	"simplebank/pkg/config"
	db "simplebank/pkg/db/sqlc"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/pressly/goose/v3"
)

func main() {

	// Get the config file path from the environment variable
	configFile := os.Getenv("CONFIG_FILE")

	// If CONFIG_FILE is not set, use a default (fallback) value
	if configFile == "" {
		configFile = "config/default.yaml" // Fallback to default.yaml if not set
	}
	
	// Load the configuration file
	config.LoadConfigs(configFile)
	cfg := config.GetConfigs()
	dbstring := fmt.Sprintf("postgresql://%v:%v@%v/%v?sslmode=%v", cfg.Postgres.UserName, url.QueryEscape(cfg.Postgres.Password), cfg.Postgres.Host, cfg.Postgres.Database, cfg.Postgres.SSLmode)
	dbConn, err := sql.Open("postgres", dbstring)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	if cfg.Postgres.Automigrate {
		// Run Goose migrations
		err = goose.Up(dbConn, "db/migrations")
		if err != nil {
			log.Fatalf("Failed to run migrations: %v\n", err)
		}
	}

	store := db.NewStore(dbConn)
	server := api.NewServer(store)
	appListenAddress := fmt.Sprintf("%v:%v", cfg.App.Host, cfg.App.Port)
	err = server.Start(appListenAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	fmt.Println("test")
}
