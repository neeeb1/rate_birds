package main

import (
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

	//err = apiCfg.PopulateBirdDB()
	//if err != nil {
	//	fmt.Printf("failed to populate birds: %s", err)
	//	return
	//}

	err = apiCfg.PopulateRatingsDB()
	if err != nil {
		fmt.Printf("failed to populate ratings: %s", err)
		return
	}

	/* 	fmt.Println("Getting top 10 birds...")
	   	topBirds, err := apiCfg.DbQueries.GetTopRatings(context.Background(), 10)
	   	if err != nil {
	   		fmt.Println("Unable to get top birds, exiting...")
	   		return
	   	}

	   	for i, b := range topBirds {
	   		birdDb, _ := apiCfg.DbQueries.GetBirdByID(context.Background(), b.BirdID)
	   		fmt.Printf("%d. %s (%d)\n", i+1, birdDb.CommonName.String, b.Rating.Int32)
	   	} */

	server.StartServer(apiCfg)
}
