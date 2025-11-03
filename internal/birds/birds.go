package birds

import (
	"context"
	"database/sql"
	"fmt"
	"math"

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
			params := database.CreateBirdParams{
				CommonName:     sql.NullString{String: b.Name, Valid: true},
				ScientificName: sql.NullString{String: b.SciName, Valid: true},
				Family:         sql.NullString{String: b.Family, Valid: true},
				Order:          sql.NullString{String: b.Order, Valid: true},
				Status:         sql.NullString{String: b.Status, Valid: true},
			}

			_, err := cfg.DbQueries.CreateBird(context.Background(), params)
			if err != nil {
				return fmt.Errorf("failed to create database entry for bird: %s", err)
			}
		}

		page++
	}

	return nil
}

func (cfg *ApiConfig) PopulateRatingsDB() error {
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

	return nil
}
