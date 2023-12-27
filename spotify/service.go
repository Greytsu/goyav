package spotify

import (
	"context"
	"golang.org/x/oauth2"

	"github.com/rs/zerolog/log"
	"github.com/zmb3/spotify/v2"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"

	"goyav/goyavUser"
)

type Service struct {
	auth *spotifyAuth.Authenticator
}

func NewService(auth *spotifyAuth.Authenticator) *Service {
	return &Service{auth: auth}
}

func (s *Service) CheckUser(tok *oauth2.Token) bool {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	// retrieve current user info
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser")
		return false
	}
	if user.ID == "" {
		return false
	}
	return true
}

func (s *Service) GetCurrentUser(tok *oauth2.Token) (goyavUser.User, error) {
	// new spotify client
	client := spotify.New(s.auth.Client(context.Background(), tok))

	// retrieve current user info
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Info().Err(err).Msg("Error while fetching for current goyavUser")
		return goyavUser.User{}, err
	}
	goyavCurrentUser := goyavUser.MapSpotifyUserToGoyavUser(user)
	return goyavCurrentUser, nil
}
