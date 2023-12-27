package goyavUser

import "time"

type User struct {
	ID          string    `bson:"user_id" json:"user_id"`
	DisplayName string    `bson:"display_name" json:"display_name"`
	SpotifyUrl  string    `bson:"spotify_url" json:"spotify_url"`
	Image       string    `bson:"image" json:"image"`
	Icon        string    `bson:"icon" json:"icon"`
	CreatedAt   time.Time `bson:"created_at,omitempty" json:"-"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty" json:"-"`
}
