package main

import (
	"github.com/cshong0618/haruka/auth/cmd/auth/subscriber"
	"github.com/cshong0618/haruka/auth/cmd/wire"
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

	authService := wire.InitAuthenticationService()
	authSubscriber := subscriber.NewAuthenticationSubscriber(authService)
	nc.QueueSubscribe("auth.create", "authapi", authSubscriber.CreateAuthentication)

	http.HandleFunc("/health", health)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
