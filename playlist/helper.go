package playlist

import (
	"errors"

	"github.com/rs/zerolog/log"

	"goyav/track"
)

var UnknownAttributeError = errors.New("unknown AttributeType")

// IsValidAttributeValue checks if the AttributeValue is within the valid bounds for the given AttributeType
func IsValidAttributeValue(attrType AttributeType, attrValue float32) (bool, error) {
	switch attrType {
	case Acousticness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Danceability:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Energy:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Instrumentalness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Key:
		return attrValue >= 0.0 && attrValue <= 11.0, nil
	case Liveness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Loudness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Mode:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Popularity:
		return attrValue >= 0.0 && attrValue <= 100.0, nil
	case Speechiness:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	case Tempo:
		return attrValue >= 0.0 && attrValue <= 300.0, nil
	case Valence:
		return attrValue >= 0.0 && attrValue <= 1.0, nil
	default:
		return false, UnknownAttributeError
	}
}

func checkAttributes(newPlaylist *Playlist) error {
	for _, attribute := range newPlaylist.TrackAttributes {
		res, err := IsValidAttributeValue(attribute.AttributeType, attribute.AttributeValue)
		if err != nil {
			log.Info().Err(err).Msg("Error while checking attributes")
			return err
		}
		if !res {
			return InvalidAttributesError
		}
	}
	return nil
}

func setupPlaylist(newPlaylist *Playlist) {
	// Setup playlist contributors
	if newPlaylist.Contributors == nil || len(newPlaylist.Contributors) == 0 {
		newPlaylist.Contributors = []string{newPlaylist.Owner}
	}

	// Setup playlist tracks
	if newPlaylist.Tracks == nil {
		newPlaylist.Tracks = map[string]track.Track{}
	}
}

func (p *Playlist) AddContributor(contributor string) {
	p.Contributors = append(p.Contributors, contributor)
}

func (p *Playlist) AddContributorToTrack(trackID string, contributor string) {
	if playlistTrack, ok := p.Tracks[trackID]; ok {
		playlistTrack.Contributors = append(playlistTrack.Contributors, contributor)
		p.Tracks[trackID] = playlistTrack
	}
}
