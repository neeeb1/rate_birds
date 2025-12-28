package birds

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RegisterEndpoints(mux *http.ServeMux, cfg *ApiConfig) {
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /api/scorematch/", cfg.handleScoreMatch)
	mux.HandleFunc("GET /api/leaderboard/", cfg.handleLoadLeaderboard)
	mux.HandleFunc("GET /api/image/", cfg.handleCachedImage)
}

func (cfg *ApiConfig) handleScoreMatch(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call to score match handler")
	rng_bird, err := cfg.DbQueries.GetRandomBirdWithImage(r.Context(), 2)
	if err != nil {
		fmt.Println(err)
	}

	leftBirdID, err := uuid.Parse(r.URL.Query().Get("leftBirdID"))
	if err != nil {
		fmt.Println(err)
	}
	leftBird, err := cfg.DbQueries.GetBirdByID(r.Context(), leftBirdID)
	if err != nil {
		fmt.Println(err)
	}

	rightBirdID, err := uuid.Parse(r.URL.Query().Get("rightBirdID"))
	if err != nil {
		fmt.Println(err)
	}
	rightBird, err := cfg.DbQueries.GetBirdByID(r.Context(), rightBirdID)
	if err != nil {
		fmt.Println(err)
	}

	winner := r.URL.Query().Get("winner")
	switch winner {
	case "left":
		//fmt.Printf("Winner: %s, Loser: %s\n", leftBird.CommonName.String, rightBird.CommonName.String)
		cfg.ScoreMatch(leftBird, rightBird)
	case "right":
		//fmt.Printf("Winner: %s, Loser: %s\n", rightBird.CommonName.String, leftBird.CommonName.String)
		cfg.ScoreMatch(rightBird, leftBird)
	}

	newLeftBird := rng_bird[0]
	newRightBird := rng_bird[1]

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	payload := fmt.Sprintf(
		`<div id="bird-wrapper" class="w-screen h-3/4 grid grid-flow-col justify-items-center">
           <div class="shadow-lg rounded-sm w-2/3 p-6 flex flex-col align-items-center bg-zinc-300" id="left-bird">
                <img class="card-image object-cover aspect-square object-contain h-2/3" src="/api/image?url=%s">
                <div class="flex flex-col text-center">
                    <p>%s</p>
                    <p><em>%s</em></p>
                    <Button class="rounded-md outline border-green-800 text-green-800 hover:border-transparent hover:bg-green-800 hover:text-white active:bg-green-900"
						hx-get="/api/scorematch/"
                        hx-trigger="click"
                        hx-target="#bird-wrapper"
                        hx-swap="outerHTML"
                        hx-vals='{"winner": "left", "leftBirdID": "%s", "rightBirdID": "%s"}'>
                        This one!
                    </Button>
                </div>
			</div>
            <div class="card-separator inline-block self-center">OR</div>
            <div class="shadow-lg rounded-sm w-2/3 p-6 flex flex-col align-items-center bg-zinc-300" id="right-bird">
                <img  class="card-image object-cover aspect-square box-content h-2/3" src="/api/image?url=%s">
                <div class="flex flex-col text-center">
                    <p>%s</p>
                    <p><em>%s</em></p>
                    <Button class="rounded-md outline border-green-800 text-green-800 hover:border-transparent hover:bg-green-800 hover:text-white active:bg-green-900"
					hx-get="/api/scorematch/"
                        hx-trigger="click"
                        hx-target="#bird-wrapper"
                        hx-swap="outerHTML"
                        hx-vals='{"winner": "right", "leftBirdID": "%s", "rightBirdID": "%s"}'>
                        This one!
                    </Button>
                </div>
            </div>
        </div>`,
		url.QueryEscape(newLeftBird.ImageUrls[0]),
		newLeftBird.CommonName.String,
		newLeftBird.ScientificName.String,
		newLeftBird.ID.String(),
		newRightBird.ID.String(),
		url.QueryEscape(newRightBird.ImageUrls[0]),
		newRightBird.CommonName.String,
		newRightBird.ScientificName.String,
		newLeftBird.ID.String(),
		newRightBird.ID.String(),
	)

	w.Write([]byte(payload))
}

func (cfg *ApiConfig) handleLoadLeaderboard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call to load leaderboard handler")

	listLength, err := strconv.Atoi(r.URL.Query().Get("listLength"))
	if err != nil {
		fmt.Println(err)
		return
	}

	topBirds, err := cfg.DbQueries.GetTopRatings(context.Background(), int32(listLength))
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")

	var builder strings.Builder
	builder.WriteString("<table>\n")

	for i, r := range topBirds {
		bird, err := cfg.DbQueries.GetBirdByID(context.Background(), r.BirdID)
		if err != nil {
			fmt.Println(err)
			return
		}

		row := fmt.Sprintf(
			`<tr>
				<td>%d.</td>
				<td>%s</td>
				<td>%d</td>
			</tr>
			`,
			i+1,
			bird.CommonName.String,
			r.Rating.Int32,
		)
		builder.WriteString(row)
	}

	builder.WriteString("</table>\n")

	payload := builder.String()

	w.Write([]byte(payload))
}

func (cfg *ApiConfig) handleCachedImage(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("url")
	if imageURL == "" {
		http.Error(w, "missing url parameter", http.StatusBadRequest)
		return
	}

	cacheURL := fmt.Sprintf("http://%s:1337/200x200,sc/%s", cfg.CacheHost, imageURL)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Get(cacheURL)
	if err != nil {
		http.Error(w, "failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		http.Error(w, "image not found", res.StatusCode)
		return
	}

	w.Header().Set("Content-Type", res.Header.Get("Content-Type"))
	w.Header().Set("Cache-Control", res.Header.Get("Cache-Control"))

	io.Copy(w, res.Body)
}
