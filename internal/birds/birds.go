package birds

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/neeeb1/rate_birds/internal/database"
)

type BirdsJson struct {
	Birds    []Bird `json:"entities"`
	Total    int    `json:"total"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}
type Bird struct {
	Images      []string `json:"images"`
	LengthMin   string   `json:"lengthMin"`
	LengthMax   string   `json:"lengthMax"`
	Name        string   `json:"name"`
	ID          int      `json:"id"`
	SciName     string   `json:"sciName"`
	Region      []string `json:"region"`
	Family      string   `json:"family"`
	Order       string   `json:"order"`
	Status      string   `json:"status"`
	WingspanMin string   `json:"wingspanMin,omitempty"`
	WingspanMax string   `json:"wingspanMax,omitempty"`
}

func (cfg *ApiConfig) PopulateBirdDB() error {
	fmt.Println("\n---* Populating birds db from Nuthatch API *---")

	intialFetch, err := cfg.GetNuthatchBirds(1, 1)
	if err != nil {
		return err
	}

	birdsToFetch := intialFetch.Total
	maxPageSize := 100
	page := 1

	for i := 0; i < birdsToFetch; i += maxPageSize {
		remaining := birdsToFetch - i
		pageSize := int(math.Min(float64(maxPageSize), float64(remaining)))

		birds, err := cfg.GetNuthatchBirds(page, pageSize)
		if err != nil {
			return err
		}

		for _, b := range birds.Birds {
			var imageArray []string
			if len(b.Images) != 0 {
				imageArray = b.Images
			}

			params := database.CreateBirdParams{
				CommonName:     sql.NullString{String: b.Name, Valid: true},
				ScientificName: sql.NullString{String: b.SciName, Valid: true},
				Family:         sql.NullString{String: b.Family, Valid: true},
				Order:          sql.NullString{String: b.Order, Valid: true},
				Status:         sql.NullString{String: b.Status, Valid: true},
				ImageUrls:      imageArray,
			}

			_, err := cfg.DbQueries.CreateBird(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to create database entry for bird: %s", err)
			}
		}
		page++
	}
	count, err := cfg.DbQueries.GetTotalBirdCount(context.Background())
	if err != nil {
		fmt.Printf("failed to count birds in db: %s\n", err)
	} else {
		fmt.Printf("Success - database contains %d bird entries\n", count)
	}

	return nil
}

func (cfg *ApiConfig) PopulateRatingsDB() error {
	fmt.Println("\n---* Populating ratings db *---")

	birds, err := cfg.DbQueries.GetAllBirds(context.Background())
	if err != nil {
		return err
	}

	for _, b := range birds {
		params := database.PopulateRatingParams{
			Matches: sql.NullInt32{Int32: 0, Valid: true},
			Rating:  sql.NullInt32{Int32: 1000, Valid: true},
			BirdID:  b.ID,
		}

		//fmt.Printf("adding %s to ratings with default values\n", b.CommonName.String)

		err = cfg.DbQueries.PopulateRating(context.Background(), params)
		if err != nil {
			return err
		}
	}

	count, err := cfg.DbQueries.GetTotalRatings(context.Background())
	if err != nil {
		fmt.Printf("failed to count ratings in db: %s\n", err)
	} else {
		fmt.Printf("Success - database contains %d rating entries\n", count)
	}
	return nil
}

func (cfg *ApiConfig) CacheImages() error {
	fmt.Println("\n*-- Starting initial image caching --*")
	startTime := time.Now()
	imageUrls, err := cfg.DbQueries.GetAllImageUrls(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get all image urls in db: %s", err)
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	var totalSuccess int
	var totalFailure int

	for _, urls := range imageUrls {
		for _, u := range urls {
			//fmt.Println(u)
			cacheUrl := fmt.Sprintf("http://%s:1337/%s", cfg.CacheHost, u)

			res, err := client.Get(cacheUrl)
			if err != nil {
				fmt.Printf("error caching image url (%s): %s\n", u, err)
				totalFailure += 1
				continue
			}

			io.Copy(io.Discard, res.Body)
			res.Body.Close()

			if res.StatusCode == http.StatusNotFound {
				fmt.Printf("response: not found\n")
				totalFailure += 1
				continue
			}

			fmt.Printf("Image cached (%s)\n", u)
			totalSuccess += 1
		}
	}

	fmt.Printf("Successfully cached %d images\n", totalSuccess)
	fmt.Printf("%d images failed to cache\n", totalFailure)
	timeElapsed := time.Since(startTime)
	fmt.Printf("Total cachine time: %s\n", timeElapsed)
	return nil
}
