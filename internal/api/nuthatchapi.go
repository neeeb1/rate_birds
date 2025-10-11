package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/neeeb1/rate_birds/internal/birds"
	"github.com/neeeb1/rate_birds/internal/database"
	"github.com/pressly/goose/v3/database"
)

type ApiConfig struct {
	NuthatcherApiKey string
	DbURL            string
	DbQueries        *database.Queries
}

func (cfg *ApiConfig) GetNuthatchBirds() (birds.BirdsJson, error) {
	fmt.Println("fetching birds from Nuthatch API...")
	var birdsJson birds.BirdsJson

	url := "https://nuthatch.lastelm.software/v2/birds"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("API-Key", cfg.NuthatcherApiKey)
	req.Header.Add("accept", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("statuscode error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &birdsJson)

	for _, b := range birdsJson.Birds {
		fmt.Println(b.Name)
		fmt.Println(b.SciName)
		fmt.Println(b.Family)
		fmt.Println()
	}

	return birdsJson, nil
}
