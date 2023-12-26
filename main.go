package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}
}

func main() {
	log.Info().Msg("hello world !")
}
