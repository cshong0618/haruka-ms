// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"github.com/cshong0618/haruka/auth/internal/eventpublisher"
	"github.com/cshong0618/haruka/auth/internal/port"
	"github.com/cshong0618/haruka/auth/internal/repository"
	mongo2 "github.com/cshong0618/haruka/auth/internal/repository/mongo"
	redis2 "github.com/cshong0618/haruka/auth/internal/repository/redis"
	"github.com/cshong0618/haruka/auth/pkg/domain/auth"
	"github.com/cshong0618/haruka/auth/pkg/usecase"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

// Injectors from wire.go:

func InitAuthenticationService() *usecase.AuthenticationService {
	client := initMongo()
	authenticationRepository := initAuthRepository(client)
	passwordHandler := port.NewPasswordHandler()
	redisClient := initRedis()
	tokenRepository := initTokenRepository(redisClient)
	conn := initNats()
	eventPublisher := eventpublisher.NewEventPublisher(conn)
	authenticationService := usecase.NewAuthenticationService(authenticationRepository, passwordHandler, tokenRepository, eventPublisher)
	return authenticationService
}

func GetNats() *nats.Conn {
	conn := initNats()
	return conn
}

// wire.go:

var natsSet = wire.NewSet(
	initNats,
)

var authServiceSet = wire.NewSet(
	initMongo,
	initRedis,
	initAuthRepository,
	natsSet, wire.Bind(new(auth.AuthenticationRepository), new(*repository.AuthenticationRepository)), initTokenRepository, wire.Bind(new(auth.TokenRepository), new(*repository.TokenRepository)), port.NewPasswordHandler, wire.Bind(new(auth.PasswordHandler), new(*port.PasswordHandler)), eventpublisher.NewEventPublisher, wire.Bind(new(auth.EventPublisher), new(*eventpublisher.EventPublisher)), usecase.NewAuthenticationService,
)

func initMongo() *mongo.Client {
	mongoURL := os.Getenv("MONGO_URL")
	return mongo2.InitMongo(mongoURL)
}

func initRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	return redis2.InitRedis(redisURL)
}

var _nats *nats.Conn

func initNats() *nats.Conn {
	if _nats != nil {
		return _nats
	}

	uri := os.Getenv("NATS_URI")
	var err error
	var nc *nats.Conn

	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}
		logrus.Infof("waiting for retry. current attempt: %d", i+1)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic(err)
	}
	logrus.Infof("Connected to %s", uri)
	_nats = nc
	return _nats
}

func initAuthRepository(mongoClient *mongo.Client) *repository.AuthenticationRepository {
	return repository.NewAuthenticationRepository(mongoClient, "auth", "auth")
}

func initTokenRepository(redisClient *redis.Client) *repository.TokenRepository {
	return repository.NewTokenRepository(redisClient)
}