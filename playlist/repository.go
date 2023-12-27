package playlist

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"

	mongoGoyav "goyav/mongo"
)

const (
	fPlaylistId = "id"
)

var PlaylistNotFoundError = errors.New("Playlist not found")

type Repository struct {
	mongoService *mongoGoyav.MongoService
	sync.Mutex
}

func NewRepository(mongoService *mongoGoyav.MongoService) *Repository {
	return &Repository{mongoService: mongoService}
}

func (repo *Repository) GetPlaylists(ctx context.Context) (*[]Playlist, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repo.mongoService.SelectionTimeout)
	defer cancel()

	// filter
	filter := bson.M{}

	var playlists []Playlist
	cursor, err := repo.mongoService.ColPlaylists().Find(timeoutCtx, filter)
	if err != nil {
		log.Info().Err(err).Msg("error while fetching playlists")
		return &[]Playlist{}, PlaylistNotFoundError
	}

	for cursor.Next(context.Background()) {
		var playlist Playlist
		if err := cursor.Decode(&playlist); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return &playlists, nil
}

func (repo *Repository) GetPlaylist(ctx context.Context, id string) (*Playlist, error) {
	timeoutCtx, cancel := context.WithTimeout(ctx, repo.mongoService.SelectionTimeout)
	defer cancel()

	// filter
	filter := bson.M{fPlaylistId: id}

	var playlist Playlist
	err := repo.mongoService.ColPlaylists().FindOne(timeoutCtx, filter).Decode(&playlist)
	if err != nil {
		log.Info().Err(err).Msg("error while fetching playlist")
		return &Playlist{}, PlaylistNotFoundError
	}
	return &playlist, nil
}

func (repo *Repository) CreatePlaylist(ctx context.Context, playlist *Playlist) (*Playlist, error) {
	//Non-thread-safety to avoid data corruption
	repo.Lock()
	defer repo.Unlock()

	timeoutCtx, cancel := context.WithTimeout(ctx, repo.mongoService.SelectionTimeout)
	defer cancel()

	for {
		//Generate UUID
		newUUID := uuid.NewV4().String()

		//Check if UUID already exists
		if _, err := repo.GetPlaylist(ctx, newUUID); err != nil {
			playlist.ID = newUUID
			playlist.CreatedAt = time.Now()
			playlist.UpdatedAt = time.Now()

			log.Info().Interface("playlist", playlist).Msg("Creating playlist")
			_, err := repo.mongoService.ColPlaylists().InsertOne(timeoutCtx, playlist)
			return playlist, err
		}
	}
}
