package artist

import "github.com/zmb3/spotify/v2"

type Artist struct {
	Name           string     `json:"name"`
	SpotifyId      spotify.ID `json:"spotify_id"`
	SpotifyUrl     string     `json:"spotify_url"`
	Genres         []string   `json:"genres"`
	ImageBigUrl    string     `json:"image_big_url"`
	ImageMediumUrl string     `json:"image_medium_url"`
	ImageSmallUrl  string     `json:"image_small_url"`
}
