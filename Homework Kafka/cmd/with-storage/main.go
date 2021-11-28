package main

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"homework-kafka/internal/http"
	"homework-kafka/internal/message_broker/kafka"
	"homework-kafka/internal/store/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main(){
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	urlExample := "postgres://goproject:goproject@localhost:5432/goproject"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}
	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer1")
	if err := broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(
		context.Background(),
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithBroker(broker),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}