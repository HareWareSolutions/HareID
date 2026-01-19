package db

import (
	"HareCRM/internal/config"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	dbPool *pgxpool.Pool
	once   sync.Once
)

func Inicialize(maxOpenConns, maxIdleConns int32, maxIdleTime time.Duration) error {

	var err error

	once.Do(func() {

		var dbConfig *pgxpool.Config

		dbConfig, err = pgxpool.ParseConfig(config.ConnectionString)
		if err != nil {
			log.Println(err)
			return
		}

		dbConfig.MaxConns = maxOpenConns // Máximo de conexões

		dbConfig.MinConns = maxIdleConns // Mínimo de conexões - Ficam 2 em standby aguardando

		dbConfig.MaxConnIdleTime = maxIdleTime // Tempo maximo que uma conexão pode ficar ociosa

		dbPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
	})

	if err != nil {
		return err
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Println(err)
		return fmt.Errorf("error db ping: %w", err)
	}

	fmt.Println("Db Conn OK")
	return nil
}

func GetPool() *pgxpool.Pool {
	return dbPool
}
