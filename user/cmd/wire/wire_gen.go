// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	mongo2 "github.com/cshong0618/haruka/user/internal/repository/mongo"
	user2 "github.com/cshong0618/haruka/user/internal/repository/user"
	"github.com/cshong0618/haruka/user/pkg/domain/usecase"
	"github.com/cshong0618/haruka/user/pkg/domain/user"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

// Injectors from wire.go:

func InitUserService() *usecase.UserService {
	client := initMongo()
	mongoRepository := initUserRepository(client)
	userService := usecase.NewUserService(mongoRepository)
	return userService
}

// wire.go:

var dbSet = wire.NewSet(
	initMongo,
	initUserRepository, wire.Bind(new(user.Repository), new(*user2.MongoRepository)), usecase.NewUserService,
)

func initMongo() *mongo.Client {
	mongoURL := os.Getenv("MONGO_URL")
	return mongo2.InitMongo(mongoURL)
}

func initUserRepository(mongoClient *mongo.Client) *user2.MongoRepository {
	return user2.NewMongoRepository(mongoClient, "user", "user")
}
