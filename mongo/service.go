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
	cUsers string = "users"

	// fields
	fUserId = "user_id"
)

type MongoService struct {
	Client           *mongo.Client
	dbName           string
	SelectionTimeout time.Duration
	colUser          *mongo.Collection
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
		colUser:          client.Database(dbName).Collection(cUsers),
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

	_, err = s.colUser.Indexes().CreateOne(timeoutCtx, mongo.IndexModel{
		Keys:    bson.D{{Key: fUserId, Value: 1}},
		Options: options.Index().SetName("user_idx"),
	})
	if err != nil {
		return err
	}

	return err
}

func (s *MongoService) ColUser() *mongo.Collection {
	return s.colUser
}
