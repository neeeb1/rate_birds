package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/neeeb1/rate_birds/internal/api"
)

func main() {
	godotenv.Load()

	apiCfg := api.ApiConfig{}

	apiCfg.NuthatcherApiKey = os.Getenv("NUTHATCHER_KEY")

	fmt.Println("apicfg loaded")

	apiCfg.GetNuthatchBirds()

	//server.StartServer()

}
