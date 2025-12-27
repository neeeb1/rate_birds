package birds

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	"github.com/neeeb1/rate_birds/internal/database"
)

const (
	// k-factor, maximum rating change
	k = 32
	d = 400
)

func (cfg *ApiConfig) ScoreMatch(winner, loser database.Bird) error {
	winnerDb, err := cfg.DbQueries.GetRatingByBirdID(context.Background(), winner.ID)
	if err != nil {
		return err
	}
	winnerRating := winnerDb.Rating.Int32

	loserDb, err := cfg.DbQueries.GetRatingByBirdID(context.Background(), loser.ID)
	if err != nil {
		return err
	}
	loserRating := loserDb.Rating.Int32

	winnerExpected, loserExpected := calculateExpected(int(winnerRating), int(loserRating))
	winnerDelta := calculateDelta(winnerExpected, 1.0)
	loserDelta := calculateDelta(loserExpected, 0.0)

	winnerNewRating := winnerRating + winnerDelta
	loserNewRating := loserRating + loserDelta

	winParams := database.UpdateRatingByBirdIDParams{
		Rating: sql.NullInt32{Int32: winnerNewRating, Valid: true},
		BirdID: winnerDb.BirdID,
	}
	_, err = cfg.DbQueries.UpdateRatingByBirdID(context.Background(), winParams)
	if err != nil {
		return err
	}

	loseParams := database.UpdateRatingByBirdIDParams{
		Rating: sql.NullInt32{Int32: loserNewRating, Valid: true},
		BirdID: loserDb.BirdID,
	}
	_, err = cfg.DbQueries.UpdateRatingByBirdID(context.Background(), loseParams)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//fmt.Printf("Updated ratings for %s and %s\n", winner.CommonName.String, loser.CommonName.String)
	fmt.Printf("(W) '%s': %d || (L) '%s': %d\n", winner.CommonName.String, winnerNewRating, loser.CommonName.String, loserNewRating)
	return nil
}

func calculateExpected(ratingA, ratingB int) (expectedA, expectedB float64) {
	qA := math.Pow10(ratingA / d)
	qB := math.Pow10(ratingB / d)

	expectedA = qA / (qA + qB)
	expectedB = qB / (qA + qB)
	return
}

func calculateDelta(expected, actual float64) int32 {
	return int32(k * (actual - expected))
}
