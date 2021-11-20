package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"project/internal/http"
	"project/internal/store/postgres"
)

func main(){
	urlExample := "postgres://goproject:goproject@localhost:5432/goproject"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	red := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	srv := http.NewServer(
		context.Background(),
		http.WithAddress(":8080"),
		http.WithStore(store),
		http.WithRedis(red),
		)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}