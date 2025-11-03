package birds

import (
	"fmt"
	"net/http"
)

func RegisterEndpoints(mux *http.ServeMux, cfg *ApiConfig) {
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /htmx/birds/random", cfg.handleGetRandomBird)
	mux.HandleFunc("POST /api/ratings/score_match", cfg.handleScoreMatch)
}

func (cfg *ApiConfig) handleGetRandomBird(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call to htmx handler")
	rng_bird, err := cfg.DbQueries.GetRandomBird(r.Context(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	type resp struct {
		Side   string `json:"side"`
		BirdID string `json:"birdID"`
	}

	var res resp
	side := r.URL.Query().Get("side")
	fmt.Println(side)
	res.Side = side

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	payload := fmt.Sprintf(
		`<div class="card" id="%s-bird">
        	<img class="card-image" src="assets/American_Barn_Owl,_Bear_River,_Utah_(9637780911).jpg">
        	<div class="card-text">
				<p>%s</p>
				<p><em>%s</em></p>
				<Button hx-get="/htmx/birds/random"
					hx-trigger="click"
					hx-target="#%s-bird"
					hx-swap="outerHTML"
					hx-vals='{"side": "%s"}'>
					This one!
				</Button>
			</div>
		</div>`,
		res.Side,
		rng_bird[0].CommonName.String,
		rng_bird[0].ScientificName.String,
		res.Side,
		res.Side)

	w.Write([]byte(payload))
}

func (cfg *ApiConfig) handleScoreMatch(w http.ResponseWriter, r *http.Request) {

}
