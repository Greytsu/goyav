package mongo

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"goyav/config"
)

const (
	// collections
	cUsers     string = "users"
	cPlaylists string = "playlists"

	// fields
	fUserId               = "user_id"
	fPlaylistId           = "ID"
	fPlaylistOwner        = "owner"
	fPlaylistContributors = "contributors"
)

type MongoService struct {
	Client           *mongo.Client
	dbName           string
	SelectionTimeout time.Duration
	colUsers         *mongo.Collection
	colPlaylists     *mongo.Collection
}

func NewService(conf *config.Config) (*MongoService, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.Mongo.URL), options.Client().SetServerSelectionTimeout(conf.Mongo.SelectionTimeout))
	if err != nil {
		return nil, err
	}

	dbName := "goyav"
	cs, err := connstring.Parse(conf.Mongo.URL)
	if cs.Database != "" {
		dbName = cs.Database
	}

	res := &MongoService{
		Client:           client,
		dbName:           dbName,
		SelectionTimeout: conf.Mongo.SelectionTimeout,
		colUsers:         client.Database(dbName).Collection(cUsers),
		colPlaylists:     client.Database(dbName).Collection(cPlaylists),
	}

	err = res.ensureIndexes(context.Background())
	if err != nil {
		return nil, err
	}

	log.Info().Msg("MongoService initialized")
	return res, nil
}

func (s *MongoService) ensureIndexes(ctx context.Context) error {

	timeoutCtx, cancel := context.WithTimeout(ctx, s.SelectionTimeout)
	defer cancel()

	var err error

	// users
	_, err = s.colUsers.Indexes().CreateOne(timeoutCtx, mongo.IndexModel{
		Keys:    bson.D{{Key: fUserId, Value: 1}},
		Options: options.Index().SetName("user_idx"),
	})
	if err != nil {
		return err
	}

	// playlists
	playlistIndexId := mongo.IndexModel{
		Keys:    bson.D{{Key: fPlaylistId, Value: 1}},
		Options: options.Index().SetName("playlist_id_idx"),
	}

	playlistIndexOwner := mongo.IndexModel{
		Keys:    bson.D{{Key: fPlaylistOwner, Value: 1}},
		Options: options.Index().SetName("playlist_owner_idx"),
	}

	playlistIndexContributors := mongo.IndexModel{
		Keys:    bson.D{{Key: fPlaylistContributors, Value: 1}},
		Options: options.Index().SetName("playlist_contributors_idx"),
	}

	_, err = s.colPlaylists.Indexes().CreateMany(context.Background(), []mongo.IndexModel{playlistIndexId, playlistIndexOwner, playlistIndexContributors})
	if err != nil {
		return err
	}

	return err
}

func (s *MongoService) ColUsers() *mongo.Collection {
	return s.colUsers
}

func (s *MongoService) ColPlaylists() *mongo.Collection {
	return s.colPlaylists
}
