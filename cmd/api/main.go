package main

import (
	"HareCRM/config"
	"HareCRM/internal/controllers"
	"HareCRM/internal/db"
	"HareCRM/internal/repository"
	"HareCRM/internal/services"
	"HareCRM/internal/validators"
	"log"
	"time"
)

/*
//Gerar secret_key

func init() {
	key := make([]byte, 64)
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}
	fmt.Println(key)

	string64 := base64.StdEncoding.EncodeToString(key)
	fmt.Println(string64)
}
*/

func main() {
	config.Load()

	dbConfig := dbConfig{
		url:          config.SUPABASE_URL,
		key:          config.SUPABASE_KEY,
		maxOpenConns: 20,
		maxIdleConns: 4,
		maxIdleTime:  2 * time.Minute,
	}

	configs := configs{
		api_port: config.PORT,
		db:       dbConfig,
	}

	if err := db.Inicialize(
		int32(dbConfig.maxOpenConns),
		int32(dbConfig.maxIdleConns),
		dbConfig.maxIdleTime,
	); err != nil {
		log.Fatalf("error initializing database: %s", err)
	}

	dbPool := db.GetPool()

	repository := repository.NewRepository(dbPool)
	validators := validators.NewValidator(repository)
	services := services.NewServices(repository, validators, dbPool)
	controllers := controllers.NewControllers(services)
	router := createRouter(controllers)

	application := application{
		config:     configs,
		repository: repository,
		services:   services,
	}

	if err := application.run(&router); err != nil {
		log.Fatal("error on application init: ", err)
	}

}
