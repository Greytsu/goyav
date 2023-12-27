package goyavUser

import "github.com/zmb3/spotify/v2"

func MapSpotifyUserToGoyavUser(user *spotify.PrivateUser) User {
	return User{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		SpotifyUrl:  user.ExternalURLs["spotify"],
		Icon:        user.Images[0].URL,
		Image:       user.Images[1].URL,
	}
}
