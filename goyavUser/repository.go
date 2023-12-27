package goyavUser

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoGoyav "goyav/mongo"
)

const (
	// fields
	FUserId      = "user_id"
	FDisplayName = "display_name"
	FSpotifyUrl  = "spotify_url"
	FImage       = "image"
	FIcon        = "icon"
)

type Repository struct {
	mongoService *mongoGoyav.MongoService
}

func NewRepository(mongoService *mongoGoyav.MongoService) *Repository {
	return &Repository{mongoService: mongoService}
}

func (repo *Repository) CreateUser(ctx context.Context, user User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, repo.mongoService.SelectionTimeout)
	defer cancel()

	//_, err := repo.mongoService.ColUser().InsertOne(timeoutCtx, user)

	_, err := repo.mongoService.ColUser().UpdateOne(timeoutCtx,
		bson.M{FUserId: user.ID},
		bson.M{
			"$set": bson.M{
				FDisplayName: user.DisplayName,
				FSpotifyUrl:  user.SpotifyUrl,
				FImage:       user.Image,
				FIcon:        user.Icon,
			},
			"$setOnInsert": bson.M{
				FUserId: user.ID,
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}
