package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mickyaujla/pack-calculator/server/internal/api"
	"github.com/mickyaujla/pack-calculator/server/internal/packing"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	server := api.NewServer(packing.DefaultSizes())
	log.Printf("listening on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, server.Routes()))
}
