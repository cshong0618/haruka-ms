//+build wireinject

package wire

import (
	mongo2 "github.com/cshong0618/haruka/post/internal/repository/mongo"
	"github.com/cshong0618/haruka/post/internal/repository/post"
	post2 "github.com/cshong0618/haruka/post/pkg/domain/post"
	"github.com/cshong0618/haruka/post/pkg/usecase"
	"github.com/google/wire"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

var wireSet = wire.NewSet(
	initNats,
)

var depSet = wire.NewSet(
	initMongo,
	initPostRepository,
	wire.Bind(new(post2.Repository), new(*post.EventSource)),
	usecase.NewPostService,
)

var _nats *nats.Conn

func initMongo() *mongo.Client {
	mongoURL := os.Getenv("MONGO_URL")
	return mongo2.InitMongo(mongoURL)
}

func initPostRepository(mongoClient *mongo.Client) *post.EventSource {
	return post.NewEventSource(mongoClient, "post", "eventCommands")
}

func InitPostService() *usecase.PostService {
	wire.Build(depSet)
	return &usecase.PostService{}
}

func initNats() *nats.Conn {
	if _nats != nil {
		return _nats
	}

	uri := os.Getenv("NATS_URI")
	nc, err := nats.Connect(uri, nats.MaxReconnects(10))

	if err != nil {
		panic(err)
	}

	_nats = nc
	return _nats
}

func GetNats() *nats.Conn {
	wire.Build(wireSet)
	return &nats.Conn{}
}
