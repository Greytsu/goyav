package goyavUser

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoGoyav "goyav/mongo"
)

const (
	// fields
	fUserId      = "user_id"
	fDisplayName = "display_name"
	fSpotifyUrl  = "spotify_url"
	fImage       = "image"
	fIcon        = "icon"
	fCreatedAt   = "created_at"
	fUpdatedAt   = "updated_at"
)

var UserNotFoundError = errors.New("User not found")

type Repository struct {
	mongoService *mongoGoyav.MongoService
}

func NewRepository(mongoService *mongoGoyav.MongoService) *Repository {
	return &Repository{mongoService: mongoService}
}

func (repo *Repository) GetUser(ctx context.Context, id string) (*User, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repo.mongoService.SelectionTimeout)
	defer cancel()

	// filter
	filter := bson.M{fUserId: id}

	var user User
	err := repo.mongoService.ColUsers().FindOne(timeoutCtx, filter).Decode(&user)
	if err != nil {
		return &User{}, UserNotFoundError
	}

	return &user, nil
}

func (repo *Repository) CreateUser(ctx context.Context, user User) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, repo.mongoService.SelectionTimeout)
	defer cancel()

	_, err := repo.mongoService.ColUsers().UpdateOne(timeoutCtx,
		bson.M{fUserId: user.ID},
		bson.M{
			"$set": bson.M{
				fDisplayName: user.DisplayName,
				fSpotifyUrl:  user.SpotifyUrl,
				fImage:       user.Image,
				fIcon:        user.Icon,
				fUpdatedAt:   time.Now(),
			},
			"$setOnInsert": bson.M{
				fUserId:    user.ID,
				fCreatedAt: time.Now(),
			},
		},
		options.Update().SetUpsert(true),
	)
	return err
}
