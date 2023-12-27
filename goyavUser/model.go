package goyavUser

type User struct {
	ID          string `json:"ID"`
	DisplayName string `json:"display_name"`
	SpotifyUrl  string `json:"spotify_url"`
	Image       string `json:"image"`
	Icon        string `json:"icon"`
}
