package spotify

import (
	"fmt"
	"time"

	"github.com/zmb3/spotify/v2"

	"goyav/track"
)

// Parses a string to Time using this layout : "2006-01-02T15:04:05.999999Z07:00"
func ParseExpiry(expiryStr string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.999999Z07:00"

	expiry, err := time.Parse(layout, expiryStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Now(), err
	}

	return expiry, nil
}

func CheckRecommendationRequest(recommendationRequest *RecommendationRequest) bool {
	if 0 > recommendationRequest.Amount || recommendationRequest.Amount > 100 {
		return false
	}

	if recommendationRequest.Timerange != spotify.ShortTermRange && recommendationRequest.Timerange != spotify.MediumTermRange && recommendationRequest.Timerange != spotify.LongTermRange {
		return false
	}

	if recommendationRequest.Type != Track && recommendationRequest.Type != Genre && recommendationRequest.Type != Artist {
		return false
	}

	return true
}

func MapTracksToSpotifyIDs(tracks map[string]track.Track) []spotify.ID {
	keys := make([]spotify.ID, 0, len(tracks))

	for key := range tracks {
		keys = append(keys, spotify.ID(key))
	}

	return keys
}
