package track

import (
	"github.com/zmb3/spotify/v2"
)

type Track struct {
	TrackName        string     `json:"track_name" bson:"track_name"`
	TrackId          spotify.ID `json:"track_id" bson:"track_id"`
	ArtistName       string     `json:"artist_name" bson:"artist_name"`
	ArtistId         spotify.ID `json:"artist_id" bson:"artist_id"`
	SpotifyAlbumName string     `json:"spotify_album_name" bson:"spotify_album_name"`
	SpotifyAlbumId   spotify.ID `json:"spotify_album_id" bson:"spotify_album_id"`
	SpotifyUrl       string     `json:"spotify_url" bson:"spotify_url"`
	ImageUrls        ImageUrls  `json:"image_urls" bson:"image_urls"`
	Contributors     []string   `json:"contributors" bson:"contributors"`
}

type ImageUrls struct {
	Big    string `json:"big" bson:"big"`
	Medium string `json:"medium" bson:"medium"`
	Small  string `json:"small" bson:"small"`
}
