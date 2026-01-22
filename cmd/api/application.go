package main

import (
	"HareCRM/internal/repository"
	"HareCRM/internal/services"
	"log"
	"net/http"
	"time"
)

type application struct {
	config     configs
	repository repository.Repository
	services   services.Services
}

type configs struct {
	api_port string
	db       dbConfig
}

type dbConfig struct {
	url          string
	key          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func (app *application) run(r *http.Handler) error {

	server := &http.Server{
		Addr:         app.config.api_port,
		Handler:      *r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("application started at port: %s", app.config.api_port)

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
