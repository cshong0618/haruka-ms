package repository

import (
	"context"
	"github.com/cshong0618/haruka/auth/pkg/domain/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type AuthenticationRepository struct {
	mongoCollection *mongo.Collection
}

func NewAuthenticationRepository(
	mongoClient *mongo.Client,
	database string,
	collection string) *AuthenticationRepository {
	db := mongoClient.Database(database)
	if db == nil {
		panic("failed to create user mongo repository. db is nil")
	}

	dbCollection := db.Collection(collection)
	if dbCollection == nil {
		panic("failed to create user mongo repository. collection is nil")
	}

	return &AuthenticationRepository{
		mongoCollection: dbCollection,
	}
}

func (a *AuthenticationRepository) CreateAuthenticationStore(ctx context.Context, userID string, handler string, password string) error {
	client := a.mongoCollection.Database().Client()
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	if err := session.StartTransaction(); err != nil {
		return err
	}

	checkHandlerQuery := bson.M{
		"handle": handler,
	}

	now := time.Now()

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		// check if handler exists
		checkHandlerResult := a.mongoCollection.FindOne(ctx, checkHandlerQuery)
		if checkHandlerResult.Err() == nil {
			return auth.HandlerAlreadyExists
		}

		dbPassword := DBPassword{
			Value:     password,
			CreatedOn: now,
		}

		dbAuthenticationStore := DBAuthenticationStore{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Handle:    handler,
			Passwords: []DBPassword{dbPassword},
			CreatedOn: now,
			UpdatedOn: now,
		}

		_, err = a.mongoCollection.InsertOne(ctx, dbAuthenticationStore)
		return err
	})

	return err
}

func (a *AuthenticationRepository) FindByHandler(ctx context.Context, handler string) (auth.UserAuthenticationStore, error) {
	findHandlerQuery := bson.M{
		"handler": handler,
	}

	result := a.mongoCollection.FindOne(ctx, findHandlerQuery)
	if err := result.Err(); err != nil {
		return auth.UserAuthenticationStore{}, err
	}

	var dbAuthenticationStore DBAuthenticationStore
	err := result.Decode(&dbAuthenticationStore)
	if err != nil {
		return auth.UserAuthenticationStore{}, err
	}

	authenticationStore := auth.UserAuthenticationStore{
		UserID:  dbAuthenticationStore.UserID,
		Handler: dbAuthenticationStore.Handle,
	}

	passwords := make([]auth.Password, len(dbAuthenticationStore.Passwords))
	for i, dbPassword := range dbAuthenticationStore.Passwords {
		passwords[i] = auth.Password{
			Value:     dbPassword.Value,
			CreatedOn: dbPassword.CreatedOn,
		}
	}

	authenticationStore.Passwords = passwords
	return authenticationStore, nil
}

type DBAuthenticationStore struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"userId"`
	Handle    string             `bson:"handle"`
	Passwords []DBPassword       `bson:"passwords"`
	CreatedOn time.Time          `bson:"createdOn"`
	UpdatedOn time.Time          `bson:"updatedOn"`
}

type DBPassword struct {
	Value     string    `bson:"value"`
	CreatedOn time.Time `bson:"createdOn"`
}
