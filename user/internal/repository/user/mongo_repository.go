package user

import (
	"context"
	userDomain "github.com/cshong0618/haruka/user/pkg/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type MongoRepository struct {
	mongoCollection *mongo.Collection
}

func NewMongoRepository(
	mongoClient *mongo.Client,
	database string,
	collection string) *MongoRepository {
	db := mongoClient.Database(database)
	if db == nil {
		panic("failed to create user mongo repository. db is nil")
	}

	dbCollection := db.Collection(collection)
	if dbCollection == nil {
		panic("failed to create user mongo repository. collection is nil")
	}

	return &MongoRepository{
		mongoCollection: dbCollection,
	}
}

func (m *MongoRepository) Create(ctx context.Context, user userDomain.User) (userDomain.User, error) {
	now := time.Now()
	dbUser := DBUser{
		ID:        primitive.NewObjectID(),
		Name:      user.Name,
		CreatedOn: now,
		UpdatedOn: now,
	}

	result, err := m.mongoCollection.InsertOne(ctx, dbUser)
	if err != nil {
		return userDomain.User{}, nil
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}

func (m *MongoRepository) FindById(ctx context.Context, ID string) (userDomain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return userDomain.User{}, err
	}

	query := bson.M{
		"_id": objectID,
	}

	result := m.mongoCollection.FindOne(ctx, query)
	if err := result.Err(); err != nil {
		return userDomain.User{}, err
	}

	var dbUser DBUser
	err = result.Decode(&dbUser)
	if err := result.Err(); err != nil {
		return userDomain.User{}, err
	}

	domainUser := userDomain.User{
		ID:   dbUser.ID.Hex(),
		Name: dbUser.Name,
	}

	return domainUser, nil
}

func (m *MongoRepository) FindAll(ctx context.Context) ([]userDomain.User, error) {
	panic("implement me")
}

type DBUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	CreatedOn time.Time          `bson:"createdOn"`
	UpdatedOn time.Time          `bson:"updatedOn"`
}
