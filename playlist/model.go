package playlist

import (
	"goyav/track"
	"time"
)

type AttributeType string

const (
	Acousticness     AttributeType = "Acousticness"
	Danceability     AttributeType = "Danceability"
	Energy           AttributeType = "Energy"
	Instrumentalness AttributeType = "Instrumentalness"
	Key              AttributeType = "Key"
	Liveness         AttributeType = "Liveness"
	Loudness         AttributeType = "Loudness"
	Mode             AttributeType = "Mode"
	Popularity       AttributeType = "Popularity"
	Speechiness      AttributeType = "Speechiness"
	Tempo            AttributeType = "Tempo"
	Valence          AttributeType = "Valence"
)

type Playlist struct {
	ID              string                 `bson:"id" json:"ID"`
	SpotifyID       string                 `bson:"spotify_id" json:"spotify_id"`
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
