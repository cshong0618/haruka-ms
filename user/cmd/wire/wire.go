//+build wireinject

package wire

import (
	mongo2 "github.com/cshong0618/haruka/user/internal/repository/mongo"
	"github.com/cshong0618/haruka/user/internal/repository/user"
	userDomain "github.com/cshong0618/haruka/user/pkg/domain/user"
	"github.com/cshong0618/haruka/user/pkg/usecase"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

var dbSet = wire.NewSet(
	initMongo,
	initUserRepository,
	wire.Bind(new(userDomain.Repository), new(*user.MongoRepository)),
	usecase.NewUserService,
	)

func initMongo() *mongo.Client {
	mongoURL := os.Getenv("MONGO_URL")
	return mongo2.InitMongo(mongoURL)
}

func initUserRepository(mongoClient *mongo.Client) *user.MongoRepository {
	return user.NewMongoRepository(mongoClient, "user", "user")
}

func InitUserService() *usecase.UserService {
	wire.Build(dbSet)
	return &usecase.UserService{}
}