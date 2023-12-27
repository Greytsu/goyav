package config

import (
	"os"
	"time"

	spotifyAuth "github.com/zmb3/spotify/v2/auth"
)

type Config struct {
	Mongo   MongoConfig
	Spotify SpotifyConfig
}

type MongoConfig struct {
	URL              string
	SelectionTimeout time.Duration
}

type SpotifyConfig struct {
	Scopes []string
}

func NewConfig() *Config {
	return &Config{
		Mongo: MongoConfig{
			URL:              os.Getenv("MONGO_URL"),
			SelectionTimeout: 10 * time.Second,
		},
		Spotify: SpotifyConfig{Scopes: []string{
			spotifyAuth.ScopePlaylistReadPrivate,
			spotifyAuth.ScopeUserReadEmail,
			spotifyAuth.ScopeUserReadEmail,
		},
		},
	}
}
