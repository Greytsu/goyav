package track

import "github.com/zmb3/spotify/v2"

type Track struct {
	Name           string     `json:"name"`
	SpotifyId      spotify.ID `json:"spotify_id"`
	SpotifyAlbumId spotify.ID `json:"spotify_album_id"`
	SpotifyUrl     string     `json:"spotify_url"`
	ImageBigUrl    string     `json:"image_big_url"`
	ImageMediumUrl string     `json:"image_medium_url"`
	ImageSmallUrl  string     `json:"image_small_url"`
}
