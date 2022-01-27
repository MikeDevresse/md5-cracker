package main

import (
	"github.com/MikeDevresse/md5-cracker/slave/internal/slave"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctrlc := make(chan os.Signal)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctrlc
		os.Exit(1)
	}()

	u := url.URL{Scheme: "ws", Host: getEnv("SLAVE_HOST", "go"), Path: getEnv("SLAVE_PATH", "/ws")}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Could not connect to master:", err)
	}
	client := slave.NewClient(c)
	client.Start()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
