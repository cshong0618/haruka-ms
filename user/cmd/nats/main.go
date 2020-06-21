package main

import (
	"github.com/cshong0618/haruka/user/cmd/nats/subscriber"
	"github.com/cshong0618/haruka/user/cmd/wire"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	uri := os.Getenv("NATS_URI")
	port := os.Getenv("PORT")
	var err error
	var nc *nats.Conn

	for i := 0; i < 10; i++ {
		nc, err = nats.Connect(uri)
		if err == nil {
			break
		}

		log.Infof("waiting for retry. current attempt: %d", i + 1)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic(err)
	}
	log.Infof("Connected to %s", uri)
	defer nc.Close()

	userService := wire.InitUserService()

	userSubscriber := subscriber.NewUserSubscriber(userService)
	nc.QueueSubscribe("user.create", "userapi", userSubscriber.CreateUser)
	nc.QueueSubscribe("user.get", "userapi", userSubscriber.GetUser)
	nc.QueueSubscribe("user.activate", "userapi", userSubscriber.ActivateUser)
	nc.QueueSubscribe("auth.creation.success", "userapi", userSubscriber.ActivateUser)
	nc.QueueSubscribe("auth.creation.failed", "userapi", userSubscriber.DeactivateUser)

	http.HandleFunc("/health", health)
	if err := http.ListenAndServe(":" + port, nil); err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}