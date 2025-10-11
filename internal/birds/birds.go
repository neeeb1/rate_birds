package birds

import (
	"context"
	"database/sql"
	"fmt"

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

func (cfg *ApiConfig) PopulateBirdDB(birds []Bird) error {
	for _, b := range birds {
		params := database.CreateBirdParams{
			CommonName:     sql.NullString{String: b.Name, Valid: true},
			ScientificName: sql.NullString{String: b.SciName, Valid: true},
			Family:         sql.NullString{String: b.Family, Valid: true},
			Order:          sql.NullString{String: b.Order, Valid: true},
			Status:         sql.NullString{String: b.Status, Valid: true},
		}

		_, err := cfg.DbQueries.CreateBird(context.Background(), params)
		if err != nil {
			return fmt.Errorf("failed to create database entry for bird: %s", b.Name)
		}
	}

	return nil
}
