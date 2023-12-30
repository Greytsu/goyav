package playlist

import (
	"time"

	"github.com/zmb3/spotify/v2"

	"goyav/track"
)

type AttributeType string

const (
	Acousticness     AttributeType = "acousticness"
	Danceability     AttributeType = "danceability"
	Energy           AttributeType = "energy"
	Instrumentalness AttributeType = "instrumentalness"
	Key              AttributeType = "key"
	Liveness         AttributeType = "liveness"
	Loudness         AttributeType = "loudness"
	Mode             AttributeType = "mode"
	Popularity       AttributeType = "popularity"
	Speechiness      AttributeType = "speechiness"
	Tempo            AttributeType = "tempo"
	Valence          AttributeType = "valence"
)

type Playlist struct {
	ID              string                 `bson:"id" json:"ID"`
	SpotifyID       spotify.ID             `bson:"spotify_id" json:"spotify_id"`
	Owner           string                 `bson:"owner" json:"owner"`
	Name            string                 `bson:"name" json:"name"`
	Tracks          map[string]track.Track `bson:"tracks" json:"tracks"`
	Contributors    []string               `bson:"contributors" json:"contributors"`
	TrackAttributes []TrackAttribute       `bson:"track_attributes" json:"track_attributes"`
	CreatedAt       time.Time              `bson:"created_at,omitempty" json:"-"`
	UpdatedAt       time.Time              `bson:"updated_at,omitempty" json:"-"`
}

type TrackAttribute struct {
	AttributeType  AttributeType `bson:"attribute_type" json:"attribute_type"`
	AttributeValue float32       `bson:"attribute_value" json:"attribute_value"`
}
