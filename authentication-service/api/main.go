package main

import (
	"authentication-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// Locally
// const webPort = "600"

// Docker
const webPort = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication")

	// Connect to the DB
	conn := connectToDB()

	//Setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Define Http Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() *sql.DB {
	// Using Docker
	dsn := os.Getenv("DSN")

	// Locally
	// dsn := "host=localhost port=5432 user=postgres password=postgres dbname=go_microservices sslmode=disable timezone=UTC connect_timeout=5"

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
