package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/neeeb1/rate_birds/internal/api"
	"github.com/neeeb1/rate_birds/internal/database"
)

func main() {
	godotenv.Load()

	apiCfg := api.ApiConfig{}

	apiCfg.NuthatcherApiKey = os.Getenv("NUTHATCH_KEY")
	apiCfg.DbURL = os.Getenv("DB_URL")

	db, err := sql.Open("postgres", apiCfg.DbURL)
	if err != nil {
		fmt.Errorf("failed to open db")
		return
	}
	apiCfg.DbQueries = database.New(db)

	fmt.Println("apicfg loaded")

	apiCfg.GetNuthatchBirds()

	//server.StartServer()

}
