package main

import (
	"fmt"
	"log"
	"net/http"
)

// Locally
// const webPort = "500"

// Docker
const webPort = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Println("Starting broker service on port: ", webPort)

	// Define Http Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
