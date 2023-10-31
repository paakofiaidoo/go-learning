package main

import (
	"database/sql"
	"flag"
	"fmt"
	"learning-go/internal/data"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	dsn  string
}
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Environment (dev|prod)")
	flag.StringVar(&cfg.dsn, "db-dsn", os.Getenv("READINGLIST_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := sql.Open("postgres", "postgres://root:root@localhost:5433/CAL?sslmode=disable")
	if err != nil {
		logger.Fatal(err)
	}
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
	}
	data.NewModels(db)
	logger.Printf("database connection pool established")
	addr := fmt.Sprintf(": %d", cfg.port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, addr)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
