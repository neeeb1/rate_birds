package server

import (
	"fmt"
	"net/http"

	"github.com/neeeb1/rate_birds/internal/birds"
)

func StartServer(cfg birds.ApiConfig) {
	mux := http.NewServeMux()
	birds.RegisterEndpoints(mux, &cfg)

	server := http.Server{}

	server.Handler = mux
	server.Addr = ":8080"

	fmt.Printf("now serving on 127.0.0.1:%s\n", server.Addr)
	server.ListenAndServe()
}
