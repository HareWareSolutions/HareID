package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ConnectionString = ""

	SUPABASE_URL = ""
	SUPABASE_KEY = ""

	PORT = ""

	SecretKey []byte
)

func Load() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	PORT = os.Getenv("API_PORT")

	SUPABASE_KEY = os.Getenv("SUPABASE_KEY")
	SUPABASE_URL = os.Getenv("SUPABASE_URL")

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
