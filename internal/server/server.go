package server

import (
	"net/http"

	"github.com/neeeb1/rate_birds/internal/birds"
)

func StartServer(cfg birds.ApiConfig) {
	mux := http.NewServeMux()
	birds.RegisterEndpoints(mux, &cfg)

	server := http.Server{}

	server.Handler = mux
	server.Addr = ":8080"

	server.ListenAndServe()
}
