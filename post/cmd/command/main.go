package main

import (
	"github.com/cshong0618/haruka/post/cmd/command/subscribers"
	"github.com/cshong0618/haruka/post/cmd/wire"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	port := os.Getenv("PORT")
	nc := wire.GetNats()
	defer nc.Close()

	encodedConn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	postService := wire.InitPostService()

	createPostCommandSubscriber := subscribers.NewCreatePostCommandSubscriber(encodedConn, postService)
	encodedConn.QueueSubscribe("post.create", "postapi", createPostCommandSubscriber.CreatePost)

	updatePostCommandSubscriber := subscribers.NewUpdatePostCommandSubscriber(encodedConn, postService)
	encodedConn.QueueSubscribe("post.update", "postapi", updatePostCommandSubscriber.UpdatePost)

	http.HandleFunc("/health", health)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r*http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}