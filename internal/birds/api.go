package birds

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func RegisterEndpoints(mux *http.ServeMux, cfg *ApiConfig) {
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /htmx/birds/random", cfg.handleGetRandomBird)
	mux.HandleFunc("POST /api/ratings/score_match", cfg.handleScoreMatch)
}

func (cfg *ApiConfig) handleGetRandomBird(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call to htmx handler")
	rng_bird, err := cfg.DbQueries.GetRandomBird(r.Context(), 2)
	if err != nil {
		fmt.Println(err)
		return
	}

	leftBirdID, err := uuid.Parse(r.URL.Query().Get("leftBirdID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	leftBird, err := cfg.DbQueries.GetBirdByID(context.Background(), leftBirdID)
	if err != nil {
		fmt.Println(err)
		return
	}

	rightBirdID, err := uuid.Parse(r.URL.Query().Get("rightBirdID"))
	if err != nil {
		fmt.Println(err)
		return
	}
	rightBird, err := cfg.DbQueries.GetBirdByID(context.Background(), rightBirdID)
	if err != nil {
		fmt.Println(err)
		return
	}

	winner := r.URL.Query().Get("winner")
	switch winner {
	case "left":
		fmt.Printf("Winner: %s, Loser: %s\n", leftBird.CommonName.String, rightBird.CommonName.String)
	case "right":
		fmt.Printf("Winner: %s, Loser: %s\n", rightBird.CommonName.String, leftBird.CommonName.String)
	}

	newLeftBird := rng_bird[0]
	newRightBird := rng_bird[1]

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	payload := fmt.Sprintf(
		`<div id="bird-wrapper">
           <div class="card" id="left-bird">
                <img class="card-image" src="assets/American_Barn_Owl,_Bear_River,_Utah_(9637780911).jpg">
                <div class="card-text">
                    <p>%s</p>
                    <p><em>%s</em></p>
                    <Button hx-get="/htmx/birds/random"
                        hx-trigger="click"
                        hx-target="#bird-wrapper"
                        hx-swap="outerHTML"
                        hx-vals='{"winner": "left", "leftBirdID": "%s", "rightBirdID": "%s"}'>
                        This one!
                    </Button>
                </div>
            </div>
            <div class="card-separator">OR</div>
            <div class="card" id="right-bird">
                <img class="card-image" src="assets/Anas_platyrhynchos_male_female_quadrat.jpg">
                <div class="card-text">
                    <p>%s</p>
                    <p><em>d%s</em></p>
                    <Button hx-get="/htmx/birds/random"
                        hx-trigger="click"
                        hx-target="#bird-wrapper"
                        hx-swap="outerHTML"
                        hx-vals='{"winner": "right", "leftBirdID": "%s", "rightBirdID": "%s"}'>
                        This one!
                    </Button>
                </div>
            </div>
        </div>`,
		newLeftBird.CommonName.String,
		newLeftBird.ScientificName.String,
		newLeftBird.ID.String(),
		newRightBird.ID.String(),
		newRightBird.CommonName.String,
		newRightBird.ScientificName.String,
		newLeftBird.ID.String(),
		newRightBird.ID.String(),
	)

	w.Write([]byte(payload))
}

func (cfg *ApiConfig) handleScoreMatch(w http.ResponseWriter, r *http.Request) {

}
