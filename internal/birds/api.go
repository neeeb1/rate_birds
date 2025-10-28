package birds

import (
	"fmt"
	"net/http"
)

func RegisterEndpoints(mux *http.ServeMux, cfg *ApiConfig) {
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /htmx/birds/random", cfg.handleGetRandomBird)

}

func (cfg *ApiConfig) handleGetRandomBird(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call to htmx handler")
	rng_bird, err := cfg.DbQueries.GetRandomBird(r.Context(), 1)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	payload := fmt.Sprintf("<p>%s</p>\n<p><em>%s</em></p>", rng_bird[0].CommonName.String, rng_bird[0].ScientificName.String)

	w.Write([]byte(payload))
}
