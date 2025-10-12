package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neeeb1/rate_birds/internal/birds"
	"github.com/neeeb1/rate_birds/internal/database"
)

func main() {
	godotenv.Load()

	apiCfg := birds.ApiConfig{}

	apiCfg.NuthatcherApiKey = os.Getenv("NUTHATCH_KEY")
	apiCfg.DbURL = os.Getenv("DB_URL")

	db, err := sql.Open("postgres", apiCfg.DbURL)
	if err != nil {
		fmt.Printf("failed to open db: %s", err)
		return
	}
	apiCfg.DbQueries = database.New(db)

	fmt.Println("apicfg loaded")

	err = apiCfg.PopulateBirdDB()
	if err != nil {
		fmt.Printf("failed to populate birds: %s", err)
		return
	}
	//server.StartServer()

}
