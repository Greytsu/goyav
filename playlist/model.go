package playlist

import "time"

type Playlist struct {
	ID           string              `bson:"id" json:"ID"`
	Owner        string              `bson:"owner" json:"owner"`
	Name         string              `bson:"name" json:"name"`
	Tracks       map[string][]string `bson:"tracks" json:"tracks"`
	Contributors []string            `bson:"contributors" json:"contributors"`
	CreatedAt    time.Time           `bson:"created_at,omitempty" json:"-"`
	UpdatedAt    time.Time           `bson:"updated_at,omitempty" json:"-"`
}
