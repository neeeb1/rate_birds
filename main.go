package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neeeb1/rate_birds/internal/birds"
	"github.com/neeeb1/rate_birds/internal/database"
	"github.com/neeeb1/rate_birds/internal/server"
)

func main() {
	if !isRunningInDockerContainer() {
		fmt.Println("Running locally, reading local .env")
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Failed to load .env: %s\n", err)
			return
		}
	} else {
		fmt.Println("Running in docker container, skipping local .env read")
	}

	apiCfg := birds.ApiConfig{}

	apiCfg.NuthatcherApiKey = os.Getenv("NUTHATCH_KEY")
	apiCfg.DbURL = os.Getenv("DB_URL")
	apiCfg.CacheHost = os.Getenv("CACHE_HOST")

	db, err := sql.Open("postgres", apiCfg.DbURL)
	if err != nil {
		fmt.Printf("failed to open db: %s", err)
		return
	}
	apiCfg.DbQueries = database.New(db)

	fmt.Println("apicfg loaded")

	count, err := apiCfg.DbQueries.GetTotalBirdCount(context.Background())
	if err != nil {
		fmt.Printf("failed to count db entries: %s", err)
		return
	}
	if count == 0 {
		err = apiCfg.PopulateBirdDB()
		if err != nil {
			fmt.Printf("failed to populate birds: %s", err)
			return
		}
	} else {
		fmt.Println("Bird db already populated - skipping initial population...")
	}

	err = apiCfg.PopulateRatingsDB()
	if err != nil {
		fmt.Printf("failed to populate ratings: %s", err)
		return
	}

	/* 	err = apiCfg.CacheImages()
	   	if err != nil {
	   		fmt.Printf("failed to cache remote images: %s", err)
	   		return
	   	} */

	server.StartServer(apiCfg)
}

func isRunningInDockerContainer() bool {
	// docker creates a .dockerenv file at the root
	// of the directory tree inside the container.
	// if this file exists then the viewer is running
	// from inside a container so return true

	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	return false
}
