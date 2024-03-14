package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// EnvMongoURI : Obtiene la cadena de conexion a MongoDB
func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("MONGOURI")
}
