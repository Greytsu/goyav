package spotify

import (
	"errors"

	"github.com/zmb3/spotify/v2"
)

var (
	InvalidRecommendationRequestError = errors.New("invalid recommendation request")
)

const (
	Track  RecommendationType = "track"
	Artist RecommendationType = "artist"
	Genre  RecommendationType = "genre"
)

type RecommendationType string

type RecommendationRequest struct {
	Amount    int                `json:"amount"`
	Type      RecommendationType `json:"type"`
	Timerange spotify.Range      `json:"timerange"`
}
