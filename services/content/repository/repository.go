package repository

import (
	"context"
	"errors"
	"time"

	"spotiftn/content/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContentRepository interface {
	CreateArtist(ctx context.Context, artist *models.Artist) (*models.Artist, error)
	UpdateArtist(ctx context.Context, artist *models.Artist) error
	GetArtistByID(ctx context.Context, id string) (*models.Artist, error)
	GetAllArtists(ctx context.Context) ([]*models.Artist, error)

	CreateAlbum(ctx context.Context, album *models.Album) (*models.Album, error)
	GetAlbumByID(ctx context.Context, id string) (*models.Album, error)
	GetAlbumsByArtist(ctx context.Context, artistID string) ([]*models.Album, error)

	CreateSong(ctx context.Context, song *models.Song) (*models.Song, error)
	GetSongsByAlbumID(ctx context.Context, albumID string) ([]*models.Song, error)
}

type MongoContentRepository struct {
	Client   *mongo.Client
	Database string
}

func NewMongoContentRepository(client *mongo.Client, dbName string) *MongoContentRepository {
	return &MongoContentRepository{
		Client:   client,
		Database: dbName,
	}
}

func (r *MongoContentRepository) CreateArtist(ctx context.Context, artist *models.Artist) (*models.Artist, error) {
	collection := r.Client.Database(r.Database).Collection("artists")
	artist.CreatedAt = time.Now()
	artist.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, artist)
	if err != nil {
		return nil, err
	}
	artist.ID = result.InsertedID.(primitive.ObjectID)
	return artist, nil
}

func (r *MongoContentRepository) UpdateArtist(ctx context.Context, artist *models.Artist) error {
	collection := r.Client.Database(r.Database).Collection("artists")

	filter := bson.M{"_id": artist.ID}
	update := bson.M{
		"$set": bson.M{
			"name":       artist.Name,
			"biography":  artist.Biography,
			"genres":     artist.Genres,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("artist not found")
	}
	return nil
}

func (r *MongoContentRepository) GetArtistByID(ctx context.Context, id string) (*models.Artist, error) {
	collection := r.Client.Database(r.Database).Collection("artists")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var artist models.Artist
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&artist)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("artist not found")
		}
		return nil, err
	}
	return &artist, nil
}

func (r *MongoContentRepository) GetAllArtists(ctx context.Context) ([]*models.Artist, error) {
	collection := r.Client.Database(r.Database).Collection("artists")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var artists []*models.Artist
	if err := cursor.All(ctx, &artists); err != nil {
		return nil, err
	}
	return artists, nil
}

func (r *MongoContentRepository) CreateAlbum(ctx context.Context, album *models.Album) (*models.Album, error) {
	collection := r.Client.Database(r.Database).Collection("albums")
	album.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, album)
	if err != nil {
		return nil, err
	}
	album.ID = result.InsertedID.(primitive.ObjectID)
	return album, nil
}

func (r *MongoContentRepository) GetAlbumByID(ctx context.Context, id string) (*models.Album, error) {
	collection := r.Client.Database(r.Database).Collection("albums")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var album models.Album
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&album)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("album not found")
		}
		return nil, err
	}
	return &album, nil
}

func (r *MongoContentRepository) GetAlbumsByArtist(ctx context.Context, artistID string) ([]*models.Album, error) {
	collection := r.Client.Database(r.Database).Collection("albums")
	objID, err := primitive.ObjectIDFromHex(artistID)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"artist_ids": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var albums []*models.Album
	if err := cursor.All(ctx, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func (r *MongoContentRepository) CreateSong(ctx context.Context, song *models.Song) (*models.Song, error) {
	collection := r.Client.Database(r.Database).Collection("songs")
	song.CreatedAt = time.Now()

	result, err := collection.InsertOne(ctx, song)
	if err != nil {
		return nil, err
	}
	song.ID = result.InsertedID.(primitive.ObjectID)
	return song, nil
}

func (r *MongoContentRepository) GetSongsByAlbumID(ctx context.Context, albumID string) ([]*models.Song, error) {
	collection := r.Client.Database(r.Database).Collection("songs")
	objID, err := primitive.ObjectIDFromHex(albumID)
	if err != nil {
		return nil, err
	}

	cursor, err := collection.Find(ctx, bson.M{"album_id": objID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var songs []*models.Song
	if err := cursor.All(ctx, &songs); err != nil {
		return nil, err
	}
	return songs, nil
}
